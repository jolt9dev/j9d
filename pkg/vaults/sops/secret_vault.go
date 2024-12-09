package sops

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/joho/godotenv"
	"github.com/jolt9dev/j9d/pkg/vaults"
	"github.com/jolt9dev/j9d/pkg/xexec"
	"github.com/jolt9dev/j9d/pkg/xfs"
)

type SopsCliSecretVault struct {
	params   SopsSecretVaultParams
	fileType string
	data     map[string]interface{}
	loaded   bool
}

type SopsAgeParams struct {
	Recipients string
	KeyFile    string
	Key        string
}

type SopsKmsParams struct {
	Uri               string
	AwsProfile        string
	EncryptionContext string
}

type SopsSecretVaultParams struct {
	File       string
	ConfigFile string
	Age        *SopsAgeParams
	KmsArns    string
	// the URI for the azure key vault and the key name
	AzureKvUri string
	// The transit engine URI for hashicorp vault
	VaultUri        string
	GcpKmsUri       string
	PgpFingerprints string
	Driver          string
	Indent          int
}

func New(params SopsSecretVaultParams) *SopsCliSecretVault {
	if params.Driver == "" {
		params.Driver = "age"
	}

	return &SopsCliSecretVault{
		params:   params,
		fileType: "dotenv",
	}
}

func init() {
	x := New(SopsSecretVaultParams{})
	var v vaults.SecretVault
	v = x

	println(v)
}

func (s *SopsCliSecretVault) LoadData(data map[string]interface{}) error {
	s.data = data
	s.loaded = true
	return nil
}

func (s *SopsCliSecretVault) GetSecretValue(key string, params *vaults.GetSecretValueParams) (string, error) {
	if !s.loaded {
		err := s.Decrypt()
		if err != nil {
			return "", err
		}
	}

	if s.fileType == "dotenv" {
		key = normalizeKey(key, s.fileType)
		if v, ok := s.data[key]; ok {
			return v.(string), nil
		}

		return "", fmt.Errorf("key not found: %s", key)
	} else {
		return "", fmt.Errorf("unsupported file type: %s", s.fileType)
	}
}

func (s *SopsCliSecretVault) ListSecretNames(params *vaults.ListSecretNamesParams) ([]string, error) {
	if !s.loaded {
		err := s.Decrypt()
		if err != nil {
			return nil, err
		}
	}

	if s.fileType == "dotenv" {
		keys := make([]string, 0, len(s.data))
		for k := range s.data {
			keys = append(keys, k)
		}

		return keys, nil
	}

	return nil, fmt.Errorf("unsupported file type: %s", s.fileType)
}

func (s *SopsCliSecretVault) BatchGetSecretValues(keys []string, params *vaults.GetSecretValueParams) (map[string]string, error) {
	values := map[string]string{}
	for _, key := range keys {
		v, err := s.GetSecretValue(key, params)
		if err != nil {
			return nil, err
		}

		values[key] = v
	}

	return values, nil
}

func (s *SopsCliSecretVault) MapSecretValues(query map[string]string, params *vaults.GetSecretValueParams) (map[string]string, error) {
	keys := make([]string, 0, len(query))
	for k := range query {
		keys = append(keys, k)
	}

	res, err := s.BatchGetSecretValues(keys, params)
	if err != nil {
		return nil, err
	}

	values := map[string]string{}
	for k, v := range query {
		if val, ok := res[k]; ok {
			values[v] = val
		}
	}

	return values, nil
}

func (s *SopsCliSecretVault) SetSecretValue(key, value string, params *vaults.SetSecretValueParams) error {
	e := s.setSecretValue(key, value)
	if e != nil {
		return e
	}

	return s.Encrypt()
}

func (s *SopsCliSecretVault) setSecretValue(key, value string) error {
	if !s.loaded {
		err := s.Decrypt()
		if err != nil {
			return err
		}
	}

	if s.fileType == "dotenv" {
		key = normalizeKey(key, s.fileType)
		s.data[key] = value
		return nil
	} else {
		return fmt.Errorf("unsupported file type: %s", s.fileType)
	}
}

func (s *SopsCliSecretVault) BatchSetSecretValues(values map[string]string, params *vaults.SetSecretValueParams) error {
	if len(values) == 0 {
		return nil
	}

	for k, v := range values {
		err := s.setSecretValue(k, v)
		if err != nil {
			return err
		}
	}

	return s.Encrypt()
}

func (s *SopsCliSecretVault) DeleteSecret(key string, params *vaults.DeleteSecretParams) error {
	if !s.loaded {
		err := s.Decrypt()
		if err != nil {
			return err
		}
	}

	if s.fileType == "dotenv" {
		key = normalizeKey(key, s.fileType)
		if _, ok := s.data[key]; ok {
			delete(s.data, key)
			err := s.Encrypt()
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("unsupported file type: %s", s.fileType)
	}

	return nil
}

func (s *SopsCliSecretVault) Decrypt() error {

	vars := map[string]string{}

	if s.fileType != "dotenv" {
		return fmt.Errorf("unsupported file type: %s", s.fileType)
	}

	switch s.params.Driver {
	case "age":
		if s.params.Age != nil {
			if s.params.Age.KeyFile != "" {
				vars["SOPS_AGE_KEY_FILE"] = s.params.Age.KeyFile
			} else if s.params.Age.Key != "" {
				vars["SOPS_AGE_KEY"] = s.params.Age.Key
			}
		}
	}

	args := []string{"decrypt"}

	if s.params.ConfigFile != "" {
		args = append(args, "--config", s.params.ConfigFile)
	}

	args = append(args, s.params.File)

	cmd := xexec.New("sops", args...)
	cmd.WithEnvMap(vars)
	cmd.WithCwd(filepath.Dir(s.params.File))

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error decrypting file: %s", out.ErrorText())
	}

	data := out.Stdout

	kv, err := godotenv.UnmarshalBytes(data)
	if err != nil {
		return err
	}

	s.data = map[string]interface{}{}
	for k, v := range kv {
		s.data[k] = v
	}

	return nil
}

func (s *SopsCliSecretVault) Encrypt() error {

	if s.fileType != "dotenv" {
		return fmt.Errorf("unsupported file type: %s", s.fileType)
	}

	vars := map[string]string{}
	args := []string{"-e"}

	if s.params.ConfigFile != "" {
		args = append(args, "--config", s.params.ConfigFile)
	}

	switch s.params.Driver {
	case "age":
		if s.params.Age != nil {
			if len(s.params.Age.Recipients) > 0 {
				args = append(args, "--age", s.params.Age.Recipients)
			}

			if s.params.Age.KeyFile != "" {
				vars["SOPS_AGE_KEY_FILE"] = s.params.Age.KeyFile
			} else if s.params.Age.Key != "" {
				vars["SOPS_AGE_KEY"] = s.params.Age.Key
			}
		}

	case "azure":
	case "azkv":
	case "azure-kv":
		if s.params.AzureKvUri != "" {
			args = append(args, "--azure-kv", s.params.AzureKvUri)
		}
	case "vault":
	case "hc-vault-transit":
		if s.params.VaultUri != "" {
			args = append(args, "--hc-vault-transit", s.params.VaultUri)
		}

	case "gcp":
	case "gcp-kms":
		if s.params.GcpKmsUri != "" {
			args = append(args, "--gcp-kms", s.params.GcpKmsUri)
		}

	case "aws":
	case "aws-kms":
	case "kms":
		if s.params.KmsArns != "" {
			args = append(args, "--kms", s.params.KmsArns)
		}

	case "pgp":
		if s.params.PgpFingerprints != "" {
			args = append(args, "--pgp", s.params.PgpFingerprints)
		}
	}

	args = append(args, "-i")
	args = append(args, s.params.File)

	kv := map[string]string{}
	for k, v := range s.data {
		kv[k] = v.(string)
	}

	str, err := godotenv.Marshal(kv)
	if err != nil {
		return err
	}

	bits := []byte(str)
	fi, err := os.Stat(s.params.File)
	var mode os.FileMode
	mode = 0644
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		mode = fi.Mode()
	}

	if xfs.Exists(s.params.File) {
		err = os.Remove(s.params.File)
	}

	// TODO: check if directory exists
	if err = os.WriteFile(s.params.File, bits, mode); err != nil {
		return err
	}

	if !xfs.Exists(s.params.File) {
		return fmt.Errorf("file not found: %s", s.params.File)
	}

	dir := filepath.Dir(s.params.File)
	for k, v := range vars {
		println(k, v)
	}

	cmd := xexec.New("sops", args...)
	cmd.WithCwd(dir)
	cmd.WithEnvMap(vars)

	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error encrypting file: %s %s %s", err.Error(), out.Text(), out.ErrorText())
	}

	return nil
}

func normalizeKey(key string, filetype string) string {
	if filetype == "dotenv" {
		sb := strings.Builder{}
		for _, c := range key {
			if c == '_' || c == '-' || c == '.' || c == '/' || c == ':' {
				sb.WriteRune('_')
				continue
			}

			if unicode.IsLetter(c) || unicode.IsDigit(c) {
				sb.WriteRune(c)
				continue
			}
		}

		return sb.String()
	}

	sb := strings.Builder{}
	for _, c := range key {
		if c == '_' || c == '-' || c == '.' || c == '/' || c == ':' {
			sb.WriteRune('.')
			continue
		}

		sb.WriteRune(c)
	}

	return sb.String()
}
