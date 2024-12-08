// pkg/vaults/sops/encrypt.go is under the Mozilla Public License Version 2.0
// https://github.com/getsops/sops/blob/main/LICENSE
//
// AFAIK there is no public api for sops encryption, so
// we can't just invoke a function to encrypt a file.
//
// The code below is remix of the sops encrypt.go and key parts of main.go
// from the sops project to enable encryption of a file.

package sops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/getsops/sops/v3"
	"github.com/getsops/sops/v3/aes"
	"github.com/getsops/sops/v3/age"
	"github.com/getsops/sops/v3/azkv"
	"github.com/getsops/sops/v3/cmd/sops/codes"
	"github.com/getsops/sops/v3/cmd/sops/common"
	"github.com/getsops/sops/v3/config"
	"github.com/getsops/sops/v3/hcvault"
	"github.com/getsops/sops/v3/keys"
	"github.com/getsops/sops/v3/keyservice"
	"github.com/getsops/sops/v3/kms"
	"github.com/getsops/sops/v3/pgp"
	"github.com/getsops/sops/v3/stores"
	"github.com/getsops/sops/v3/version"
	"github.com/mitchellh/go-wordwrap"
)

type encryptConfig struct {
	UnencryptedSuffix       string
	EncryptedSuffix         string
	UnencryptedRegex        string
	EncryptedRegex          string
	UnencryptedCommentRegex string
	EncryptedCommentRegex   string
	MACOnlyEncrypted        bool
	KeyGroups               []sops.KeyGroup
	GroupThreshold          int
}

type encryptOpts struct {
	Cipher      sops.Cipher
	InputStore  sops.Store
	OutputStore sops.Store
	InputPath   string
	KeyServices []keyservice.KeyServiceClient
	encryptConfig
}

type fileAlreadyEncryptedError struct{}

func (err *fileAlreadyEncryptedError) Error() string {
	return "File already encrypted"
}

func shamirThreshold(configPath string, file string) (int, error) {
	conf, err := loadConfig(configPath, file, nil)
	if conf == nil {
		// This takes care of the following two case:
		// 1. No config was provided, or contains no creation rules. Err will be nil and ShamirThreshold will be the default value of 0.
		// 2. We did find a config file, but failed to load it. In that case the calling function will print the error and exit.
		return 0, err
	}
	return conf.ShamirThreshold, nil
}

func loadConfig(configPath string, file string, kmsEncryptionContext map[string]*string) (*config.Config, error) {
	var err error
	if configPath == "" {
		// Ignore config not found errors returned from FindConfigFile since the config file is not mandatory
		configPath, err = config.FindConfigFile(".")
		if err != nil {
			// If we can't find a config file, but we were not explicitly requested to, assume it does not exist
			return nil, nil
		}
	}
	conf, err := config.LoadCreationRuleForFile(configPath, file, kmsEncryptionContext)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func keyGroups(configPath string, file string, params SopsEncryptParams) ([]sops.KeyGroup, error) {
	var kmsKeys []keys.MasterKey
	var pgpKeys []keys.MasterKey
	//var cloudKmsKeys []keys.MasterKey
	var azkvKeys []keys.MasterKey
	var hcVaultMkKeys []keys.MasterKey
	var ageMasterKeys []keys.MasterKey
	var specific = false

	if params.Age != nil {
		ageKeys, err := age.MasterKeysFromRecipients(strings.Join(params.Age.Recipients, ","))
		if err != nil {
			return nil, err
		}
		specific = true
		for _, k := range ageKeys {
			ageMasterKeys = append(ageMasterKeys, k)
			specific = true
		}
	}

	if params.VaultUri != "" {
		specific = true
		hcVaultKeys, err := hcvault.NewMasterKeysFromURIs(params.AzureKvUri)
		if err != nil {
			return nil, err
		}
		for _, k := range hcVaultKeys {
			hcVaultMkKeys = append(hcVaultMkKeys, k)
		}
	}

	if params.Kms != nil {
		specific = true
		kmsEncryptionContext := kms.ParseKMSContext(params.Kms.EncryptionContext)
		if params.Kms.EncryptionContext != "" && kmsEncryptionContext == nil {
			return nil, common.NewExitError("Invalid KMS encryption context format", codes.ErrorInvalidKMSEncryptionContextFormat)
		}

		for _, k := range kms.MasterKeysFromArnString(params.Kms.Uri, kmsEncryptionContext, params.Kms.AwsProfile) {
			kmsKeys = append(kmsKeys, k)
		}
	}

	if params.AzureKvUri != "" {
		specific = true
		azureKeys, err := azkv.MasterKeysFromURLs(params.AzureKvUri)
		if err != nil {
			return nil, err
		}
		for _, k := range azureKeys {
			azkvKeys = append(azkvKeys, k)
		}
	}

	if params.PgpKey != "" {
		specific = true
		for _, k := range pgp.MasterKeysFromFingerprintString(params.PgpKey) {
			pgpKeys = append(pgpKeys, k)
		}
	}

	if specific {
		var group sops.KeyGroup
		group = append(group, kmsKeys...)
		// group = append(group, cloudKmsKeys...)
		group = append(group, azkvKeys...)
		group = append(group, pgpKeys...)
		group = append(group, hcVaultMkKeys...)
		group = append(group, ageMasterKeys...)

		return []sops.KeyGroup{group}, nil
	}

	conf, err := loadConfig(configPath, file, nil)
	// config file might just not be supplied, without any error
	if conf == nil {
		errMsg := "config file not found, or has no creation rules, and no keys provided through command line options"
		if err != nil {
			errMsg = fmt.Sprintf("%s: %s", errMsg, err)
		}
		return nil, fmt.Errorf(errMsg)
	}
	return conf.KeyGroups, err
}

func (err *fileAlreadyEncryptedError) UserError() string {
	message := "The file you have provided contains a top-level entry called " +
		"'" + stores.SopsMetadataKey + "', or for flat file formats top-level entries starting with " +
		"'" + stores.SopsMetadataKey + "_'. This is generally due to the file already being encrypted. " +
		"SOPS uses a top-level entry called '" + stores.SopsMetadataKey + "' to store the metadata " +
		"required to decrypt the file. For this reason, SOPS can not " +
		"encrypt files that already contain such an entry.\n\n" +
		"If this is an unencrypted file, rename the '" + stores.SopsMetadataKey + "' entry.\n\n" +
		"If this is an encrypted file and you want to edit it, use the " +
		"editor mode, for example: `sops my_file.yaml`"
	return wordwrap.WrapString(message, 75)
}

func ensureNoMetadata(opts encryptOpts, branch sops.TreeBranch) error {
	if opts.OutputStore.HasSopsTopLevelKey(branch) {
		return &fileAlreadyEncryptedError{}
	}
	return nil
}

func metadataFromEncryptionConfig(config encryptConfig) sops.Metadata {
	return sops.Metadata{
		KeyGroups:               config.KeyGroups,
		UnencryptedSuffix:       config.UnencryptedSuffix,
		EncryptedSuffix:         config.EncryptedSuffix,
		UnencryptedRegex:        config.UnencryptedRegex,
		EncryptedRegex:          config.EncryptedRegex,
		UnencryptedCommentRegex: config.UnencryptedCommentRegex,
		EncryptedCommentRegex:   config.EncryptedCommentRegex,
		MACOnlyEncrypted:        config.MACOnlyEncrypted,
		Version:                 version.Version,
		ShamirThreshold:         config.GroupThreshold,
	}
}

func loadStoresConfig(configPath string) (*config.StoresConfig, error) {
	if configPath == "" {
		// Ignore config not found errors returned from FindConfigFile since the config file is not mandatory
		foundPath, err := config.FindConfigFile(".")
		if err != nil {
			return config.NewStoresConfig(), nil
		}
		configPath = foundPath
	}
	return config.LoadStoresConfig(configPath)
}

func inputStore(configPath string, path string, fileType string) (common.Store, error) {
	storesConf, err := loadStoresConfig(configPath)
	if err != nil {
		return nil, err
	}
	return common.DefaultStoreForPathOrFormat(storesConf, path, fileType), nil
}

func outputStore(configPath string, path string, fileType string, indent int) (common.Store, error) {
	storesConf, err := loadStoresConfig(configPath)
	if err != nil {
		return nil, err
	}
	if indent > 0 {
		storesConf.YAML.Indent = indent
		storesConf.JSON.Indent = indent
		storesConf.JSONBinary.Indent = indent
	}

	return common.DefaultStoreForPathOrFormat(storesConf, path, fileType), nil
}

func decryptionOrder(decryptionOrder string) ([]string, error) {
	if decryptionOrder == "" {
		return sops.DefaultDecryptionOrder, nil
	}
	orderList := strings.Split(decryptionOrder, ",")
	unique := make(map[string]struct{})
	for _, v := range orderList {
		if _, ok := unique[v]; ok {
			return nil, common.NewExitError(fmt.Sprintf("Duplicate decryption key type: %s", v), codes.DuplicateDecryptionKeyType)
		}
		unique[v] = struct{}{}
	}
	return orderList, nil
}

func getEncryptConfig(configPath string, fileName string, params SopsEncryptParams) (encryptConfig, error) {
	unencryptedSuffix := ""
	encryptedSuffix := ""
	encryptedRegex := ""
	unencryptedRegex := ""
	encryptedCommentRegex := ""
	unencryptedCommentRegex := ""
	macOnlyEncrypted := false
	conf, err := loadConfig(configPath, fileName, nil)
	if err != nil {
		return encryptConfig{}, err
	}
	if conf != nil {
		// command line options have precedence
		if unencryptedSuffix == "" {
			unencryptedSuffix = conf.UnencryptedSuffix
		}
		if encryptedSuffix == "" {
			encryptedSuffix = conf.EncryptedSuffix
		}
		if encryptedRegex == "" {
			encryptedRegex = conf.EncryptedRegex
		}
		if unencryptedRegex == "" {
			unencryptedRegex = conf.UnencryptedRegex
		}
		if encryptedCommentRegex == "" {
			encryptedCommentRegex = conf.EncryptedCommentRegex
		}
		if unencryptedCommentRegex == "" {
			unencryptedCommentRegex = conf.UnencryptedCommentRegex
		}
		if !macOnlyEncrypted {
			macOnlyEncrypted = conf.MACOnlyEncrypted
		}
	}

	cryptRuleCount := 0
	if unencryptedSuffix != "" {
		cryptRuleCount++
	}
	if encryptedSuffix != "" {
		cryptRuleCount++
	}
	if encryptedRegex != "" {
		cryptRuleCount++
	}
	if unencryptedRegex != "" {
		cryptRuleCount++
	}
	if encryptedCommentRegex != "" {
		cryptRuleCount++
	}
	if unencryptedCommentRegex != "" {
		cryptRuleCount++
	}

	if cryptRuleCount > 1 {
		return encryptConfig{}, common.NewExitError("Error: cannot use more than one of encrypted_suffix, unencrypted_suffix, encrypted_regex, unencrypted_regex, encrypted_comment_regex, or unencrypted_comment_regex in the same file", codes.ErrorConflictingParameters)
	}

	// only supply the default UnencryptedSuffix when EncryptedSuffix, EncryptedRegex, and others are not provided
	if cryptRuleCount == 0 {
		unencryptedSuffix = sops.DefaultUnencryptedSuffix
	}

	var groups []sops.KeyGroup
	groups, err = keyGroups(configPath, fileName, params)
	if err != nil {
		return encryptConfig{}, err
	}

	var threshold int
	threshold, err = shamirThreshold(configPath, fileName)
	if err != nil {
		return encryptConfig{}, err
	}

	return encryptConfig{
		UnencryptedSuffix:       unencryptedSuffix,
		EncryptedSuffix:         encryptedSuffix,
		UnencryptedRegex:        unencryptedRegex,
		EncryptedRegex:          encryptedRegex,
		UnencryptedCommentRegex: unencryptedCommentRegex,
		EncryptedCommentRegex:   encryptedCommentRegex,
		MACOnlyEncrypted:        macOnlyEncrypted,
		KeyGroups:               groups,
		GroupThreshold:          threshold,
	}, nil
}

type SopsEncryptParams struct {
	ConfigPath string
	File       string
	FileType   string
	Indent     int
	Age        *SopsAgeParams
	Kms        *SopsKmsParams
	AzureKvUri string
	VaultUri   string
	PgpKey     string
}

func encryptOutput(params SopsEncryptParams) ([]byte, error) {
	configPath := params.ConfigPath
	file := params.File
	fileType := params.FileType
	indent := params.Indent
	inputStore, err := inputStore(configPath, file, fileType)
	if err != nil {
		return nil, err
	}
	outputStore, err := outputStore(configPath, file, fileType, indent)
	if err != nil {
		return nil, err
	}

	svcs := []keyservice.KeyServiceClient{}
	svcs = append(svcs, keyservice.NewLocalClient())

	var output []byte

	encConfig, err := getEncryptConfig(configPath, file, params)
	if err != nil {
		return nil, err
	}
	output, err = encrypt(encryptOpts{
		OutputStore:   outputStore,
		InputStore:    inputStore,
		InputPath:     file,
		Cipher:        aes.NewCipher(),
		KeyServices:   svcs,
		encryptConfig: encConfig,
	})

	return output, err
}

func encrypt(opts encryptOpts) (encryptedFile []byte, err error) {
	// Load the file
	fileBytes, err := os.ReadFile(opts.InputPath)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error reading file: %s", err), codes.CouldNotReadInputFile)
	}
	branches, err := opts.InputStore.LoadPlainFile(fileBytes)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Error unmarshalling file: %s", err), codes.CouldNotReadInputFile)
	}
	if len(branches) < 1 {
		return nil, common.NewExitError("File cannot be completely empty, it must contain at least one document", codes.NeedAtLeastOneDocument)
	}
	if err := ensureNoMetadata(opts, branches[0]); err != nil {
		return nil, common.NewExitError(err, codes.FileAlreadyEncrypted)
	}
	path, err := filepath.Abs(opts.InputPath)
	if err != nil {
		return nil, err
	}
	tree := sops.Tree{
		Branches: branches,
		Metadata: metadataFromEncryptionConfig(opts.encryptConfig),
		FilePath: path,
	}
	dataKey, errs := tree.GenerateDataKeyWithKeyServices(opts.KeyServices)
	if len(errs) > 0 {
		err = fmt.Errorf("Could not generate data key: %s", errs)
		return nil, err
	}

	err = common.EncryptTree(common.EncryptTreeOpts{
		DataKey: dataKey,
		Tree:    &tree,
		Cipher:  opts.Cipher,
	})
	if err != nil {
		return nil, err
	}

	encryptedFile, err = opts.OutputStore.EmitEncryptedFile(tree)
	if err != nil {
		return nil, common.NewExitError(fmt.Sprintf("Could not marshal tree: %s", err), codes.ErrorDumpingTree)
	}
	return
}
