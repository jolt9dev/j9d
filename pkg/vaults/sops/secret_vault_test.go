package sops_test

import (
	"os"
	"testing"

	"github.com/jolt9dev/j9d/pkg/vaults/sops"
	"github.com/jolt9dev/j9d/pkg/xfs"
	fs "github.com/jolt9dev/j9d/pkg/xfs"
	"github.com/stretchr/testify/assert"
)

func TestSopsSecretVault(t *testing.T) {

	defer func() {
		if fs.Exists("./.env") {
			os.Remove("./.env")
		}
	}()

	publicKey := "age1690jcga9k3976xdldnk7wpypdcpryq4afmt4st3zltvxand3z97qzdewun"
	privateKey := "AGE-SECRET-KEY-1NVF2YM9LNT8ZVTZQEAC0F374Z0EJ2S57W3T0XHTZY3GMQC62NLZQJV33FU"

	// os.Setenv("SOPS_AGE_RECIPIENTS", publicKey)
	// os.Setenv("SOPS_AGE_KEY", privateKey)

	cwd, err := fs.Cwd()
	if err != nil {
		t.Fatal(err)
	}

	f, err := xfs.Resolve("./.env", cwd)
	if err != nil {
		t.Fatal(err)
	}

	vault := sops.New(sops.SopsSecretVaultParams{
		File:   f,
		Driver: "age",
		Indent: 0,
		Age: &sops.SopsAgeParams{
			Recipients: publicKey,
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
