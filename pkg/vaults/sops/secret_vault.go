package sops

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/getsops/sops/v3/decrypt"
	"github.com/joho/godotenv"
	"github.com/jolt9dev/jolt9/pkg/vaults"
)

type SopsSecretVault struct {
	params   SopsSecretVaultParams
	fileType string
	data     map[string]interface{}
	loaded   bool
}

type SopsSecretVaultParams struct {
	File         string
	ConfileFile  string
	Age          *SopsAgeParams
	Kms          *SopsKmsParams
	AzureKvUri   string
	VaultUri     string
	PgpPublicKey string
	Driver       string
	Indent       int
}

type SopsAgeParams struct {
	Recipients []string
	Key        string
}

type SopsKmsParams struct {
	Uri               string
	AwsProfile        string
	EncryptionContext string
}

func New(params SopsSecretVaultParams) *SopsSecretVault {
	if params.Driver == "" {
		params.Driver = "age"
	}

	return &SopsSecretVault{
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

func (s *SopsSecretVault) LoadData(data map[string]interface{}) error {
	s.data = data
	s.loaded = true
	return nil
}

func (s *SopsSecretVault) GetSecretValue(key string, params *vaults.GetSecretValueParams) (string, error) {
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

func (s *SopsSecretVault) ListSecretNames(params *vaults.ListSecretNamesParams) ([]string, error) {
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

func (s *SopsSecretVault) BatchGetSecretValues(keys []string, params *vaults.GetSecretValueParams) (map[string]string, error) {
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

func (s *SopsSecretVault) MapSecretValues(query map[string]string, params *vaults.GetSecretValueParams) (map[string]string, error) {
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

func (s *SopsSecretVault) SetSecretValue(key, value string, params *vaults.SetSecretValueParams) error {
	e := s.setSecretValue(key, value)
	if e != nil {
		return e
	}

	return s.Encrypt()
}

func (s *SopsSecretVault) setSecretValue(key, value string) error {
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

func (s *SopsSecretVault) BatchSetSecretValues(values map[string]string, params *vaults.SetSecretValueParams) error {
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

func (s *SopsSecretVault) DeleteSecret(key string, params *vaults.DeleteSecretParams) error {
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

func (s *SopsSecretVault) Encrypt() error {

	if s.fileType != "dotenv" {
		return fmt.Errorf("unsupported file type: %s", s.fileType)
	}

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

	// TODO: check if directory exists
	if err = os.WriteFile(s.params.File, bits, mode); err != nil {
		return err
	}

	bytes, err := encryptOutput(SopsEncryptParams{
		File:       s.params.File,
		FileType:   s.fileType,
		Indent:     0,
		ConfigPath: s.params.ConfileFile,
		Age:        s.params.Age,
		Kms:        s.params.Kms,
		AzureKvUri: s.params.AzureKvUri,
		VaultUri:   s.params.VaultUri,
		PgpKey:     s.params.PgpPublicKey,
	})

	if err != nil {
		return err
	}

	return os.WriteFile(s.params.File, bytes, mode)
}

func (s *SopsSecretVault) Decrypt() error {

	data, err := decrypt.File(s.params.File, s.fileType)
	if err != nil {
		return err
	}

	if s.fileType != "dotenv" {
		return fmt.Errorf("unsupported file type: %s", s.fileType)
	}

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
