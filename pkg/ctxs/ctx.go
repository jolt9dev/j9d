package ctxs

import (
	"fmt"
	"net/url"
	"path/filepath"

	"github.com/jolt9dev/j9d/pkg/env"
	"github.com/jolt9dev/j9d/pkg/types"
	"github.com/jolt9dev/j9d/pkg/vaults"
	"github.com/jolt9dev/j9d/pkg/vaults/sops"
	fs "github.com/jolt9dev/j9d/pkg/xfs"
	"github.com/m1/go-generate-password/generator"
	"gopkg.in/yaml.v3"
)

type ExecContext struct {
	Env     map[string]string
	Secrets map[string]string
	Jolt9   *types.Jolt9
	Cwd     string
}

func Load(file string) (*ExecContext, error) {
	if !fs.Exists(file) {
		return nil, fmt.Errorf("file %s not found", file)
	}

	workingDir := filepath.Dir(file)

	vaults := make(map[string]vaults.SecretVault)
	vars := make(map[string]string)
	secrets := make(map[string]string)

	bytes, err := fs.ReadFile(file)
	if err != nil {
		return nil, err
	}

	jolt9 := &types.Jolt9{}
	err = yaml.Unmarshal(bytes, jolt9)
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(file)

	jolt9, err = jolt9.ResolveInheritence(dir)
	if err != nil {
		return nil, err
	}

	if len(jolt9.Secrets) > 0 && len(jolt9.Vaults) == 0 {
		return nil, fmt.Errorf("secrets found but no vault")
	}

	if len(jolt9.Vaults) > 0 {
		for _, vault := range jolt9.Vaults {

			u, err := url.Parse(vault.Uri)

			if err != nil {
				return nil, err
			}

			switch u.Scheme {
			case "sops":
				v, err := loadSopsVault(&vault, workingDir)
				if err != nil {
					return nil, err
				}

				vaults[vault.Name] = v

			default:
				return nil, fmt.Errorf("unsupported vault scheme %s", u.Scheme)
			}
		}
	}

	vaultCount := len(vaults)

	for _, s := range jolt9.Secrets {
		if s.Key == "" {
			s.Key = s.Name
		}

		secretValue := ""

		if s.Vault == "" {
			for _, vt := range vaults {
				v, err := vt.GetSecretValue(s.Key, nil)
				if err == nil && v != "" {

					secretValue = v
					vars[s.Name] = secretValue
					break
				}
			}
		} else {
			vt, ok := vaults[s.Vault]
			if ok {
				v, err := vt.GetSecretValue(s.Key, nil)
				if err == nil && v != "" {

					secretValue = v
					vars[s.Name] = secretValue
				}
			}
		}

		if secretValue == "" {

			if s.Gen {
				if s.Vault == "" && vaultCount > 1 {
					return nil, fmt.Errorf("no vault specified for generated secret %s", s.Name)
				}

				if s.Vault != "" {
					return nil, fmt.Errorf("no vaults found for generated secret %s", s.Name)
				}

				if s.Size == 0 {
					s.Size = 16
				}

				chars := ""
				if s.Digits {
					chars += "0123456789"
				}

				if s.Lower {
					chars += "abcdefghijklmnopqrstuvwxyz"
				}

				if s.Upper {
					chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
				}

				if s.Special != "" {
					chars += s.Special
				}

				gen, err := generator.New(&generator.Config{
					Length:       uint(s.Size),
					CharacterSet: chars,
				})

				if err != nil {
					return nil, err
				}

				v2, err := gen.Generate()
				if err != nil {
					return nil, err
				}

				secretValue = *v2

				if vaultCount == 1 {
					for _, vt := range vaults {
						err = vt.SetSecretValue(s.Key, secretValue, nil)
						if err != nil {
							return nil, err
						}

						break
					}
				} else {
					vt, ok := vaults[s.Vault]
					if ok {
						err = vt.SetSecretValue(s.Key, secretValue, nil)
						if err != nil {
							return nil, err
						}
					}
				}
			} else {
				return nil, fmt.Errorf("secret %s not found", s.Name)
			}
		}

		secrets[s.Name] = secretValue
		env.Set(s.Name, secretValue)
		vars[s.Name] = secretValue
	}

	if len(jolt9.Env) > 0 {
		for k, v := range jolt9.Env {
			n, err := env.Expand(v, nil)
			if err != nil {
				return nil, err
			}

			env.Set(k, n)
			vars[k] = n
		}
	}

	return &ExecContext{
		Env:     vars,
		Secrets: secrets,
		Jolt9:   jolt9,
		Cwd:     workingDir,
	}, nil
}

func loadSopsVault(vault *types.Vault, cwd string) (*sops.SopsCliSecretVault, error) {
	u, err := url.Parse(vault.Uri)
	if err != nil {
		return nil, err
	}

	sopsFile := u.Path
	if u.Host == "." || u.Host == ".." {
		sopsFile = u.Host + sopsFile
	}

	sopsFile, err = fs.Resolve(sopsFile, cwd)
	if err != nil {
		return nil, err
	}
	println(sopsFile)
	recipients := u.Query().Get("age-recipients")
	sopsKeyFile := u.Query().Get("age-key-file")
	configFile := u.Query().Get("config")

	if recipients == "" {
		v, ok := vault.With["age-recipients"]
		if ok && v != nil {
			recipients = v.(string)
		}
	}

	if recipients == "" {
		recipients = env.Get("SOPS_AGE_RECIPIENTS")
	}

	if sopsKeyFile == "" {
		v, ok := vault.With["age-key-file"]
		if ok && v != nil {
			sopsKeyFile = v.(string)
		}
	}

	if sopsFile == "" {
		v, ok := vault.With["file"]
		if ok && v != nil {
			sopsFile = v.(string)
		}
	}

	if configFile == "" {
		v, ok := vault.With["config"]
		if ok && v != nil {
			configFile = v.(string)
		}
	}

	if sopsFile == "" {
		return nil, fmt.Errorf("sops file not found")
	}

	if sopsKeyFile != "" {
		n, err := fs.Resolve(sopsKeyFile, cwd)
		if err != nil {
			return nil, err
		}

		env.Set("SOPS_AGE_KEY_FILE", n)
	}

	params := &sops.SopsSecretVaultParams{
		File:       sopsFile,
		ConfigFile: configFile,
	}

	if recipients != "" {
		params.Driver = "age"
		params.Age = &sops.SopsAgeParams{
			Recipients: recipients,
			KeyFile:    sopsKeyFile,
		}
	}

	return sops.New(*params), nil
}
