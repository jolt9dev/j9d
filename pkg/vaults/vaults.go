package vaults

import "context"

type OperationParams struct {
	Context context.Context
}

type GetSecretValueParams struct {
	OperationParams
	Version string
}

type BatchGetSecretValuesParams struct {
	OperationParams
}

type SetSecretValueParams struct {
	OperationParams
}

type BatchSetSecretValuesParams struct {
	OperationParams
}

type DeleteSecretParams struct {
	OperationParams
}

type ListSecretNamesParams struct {
	OperationParams
}

type SecretVault interface {
	BatchGetSecretValues(keys []string, params *GetSecretValueParams) (map[string]string, error)

	BatchSetSecretValues(values map[string]string, params *SetSecretValueParams) error

	GetSecretValue(key string, params *GetSecretValueParams) (string, error)

	ListSecretNames(params *ListSecretNamesParams) ([]string, error)

	MapSecretValues(keys map[string]string, params *GetSecretValueParams) (map[string]string, error)

	SetSecretValue(key, value string, params *SetSecretValueParams) error

	DeleteSecret(key string, params *DeleteSecretParams) error
}
