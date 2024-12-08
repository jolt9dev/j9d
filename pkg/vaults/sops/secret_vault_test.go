package sops_test

import (
	"os"
	"testing"

	"filippo.io/age"
	fs "github.com/jolt9dev/j9d/pkg/xfs"
	"github.com/jolt9dev/jolt9/pkg/vaults/sops"
	"github.com/stretchr/testify/assert"
)

func TestSopsSecretVault(t *testing.T) {

	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if fs.Exists("./.env") {
			err = os.Remove("./.env")
		}
	}()

	publicKey := id.Recipient().String()
	privateKey := id.String()

	// os.Setenv("SOPS_AGE_RECIPIENTS", publicKey)
	// os.Setenv("SOPS_AGE_KEY", privateKey)

	vault := sops.New(sops.SopsSecretVaultParams{
		File:        "./.env",
		ConfileFile: "",
		Driver:      "age",
		Indent:      0,
		Age: &sops.SopsAgeParams{
			Recipients: []string{publicKey},
			Key:        privateKey,
		},
	})

	data := map[string]interface{}{
		"VAR1": "VALUE1",
		"VAR2": "VALUE2",
	}

	vault.LoadData(data)
	err = vault.Encrypt()
	if err != nil {
		t.Fatal(err)
	}

	secret, err := vault.GetSecretValue("VAR1", nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "VALUE1", secret)

	err = vault.SetSecretValue("NEW_VAR", "NEW_VALUE", nil)
	if err != nil {
		t.Fatal(err)
	}

	secret, err = vault.GetSecretValue("NEW_VAR", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = vault.DeleteSecret("NEW_VAR", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = vault.GetSecretValue("NEW_VAR", nil)
	assert.NotNil(t, err)
}
