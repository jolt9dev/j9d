package ssh

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"

	"github.com/kevinburke/ssh_config"
)

func FindConfig(alias string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	sshConfigFile := filepath.Join(home, ".ssh", "ssh_config")

	if _, err := os.Stat(sshConfigFile); os.IsNotExist(err) {
		return errors.New("ssh config file does not exist")
	}

	sshConfigData, err := os.ReadFile(sshConfigFile)
	if err != nil {
		return err
	}

	sshConfig, err := ssh_config.DecodeBytes(sshConfigData)
	if err != nil {
		return err
	}

	if sshConfig.Hosts == nil {
		return errors.New("ssh config file does not contain any hosts")
	}

	hostFound := false

	port := 22
	hostname := alias
	username := ""
	identityFiles := []string{}
	passwordAuth := false

	for _, host := range sshConfig.Hosts {
		if host.Matches(alias) {
			hostFound = true
			for _, n := range host.Nodes {
				if kv, ok := n.(*ssh_config.KV); ok {
					switch kv.Key {
					case "HostName":
						hostname = kv.Value
					case "User":
						username = kv.Value
					case "Port":
						portValue := kv.Value
						port, err = strconv.Atoi(portValue)
						if err != nil {
							return err
						}
					case "IdentityFile":
						identityFiles = append(identityFiles, kv.Value)
					case "PasswordAuthentication":
						passwordAuth = kv.Value == "yes"
					}
				}
				break
			}
		}
	}

	println(port)
	println(hostname)
	println(username)
	println(identityFiles)
	println(passwordAuth)

	if !hostFound {
		return errors.New("host not found in ssh config file")
	}

	return nil
}
