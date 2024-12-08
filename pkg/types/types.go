package types

import (
	"encoding/json"
	"errors"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/jolt9dev/j9d/pkg/consts"
	"github.com/jolt9dev/j9d/pkg/paths"
	fs "github.com/jolt9dev/j9d/pkg/xfs"
	"gopkg.in/yaml.v3"
)

var (
	globalConfigFile *GlobalConfigFile
)

func GetGlobalConfig() (*GlobalConfig, error) {
	if globalConfigFile != nil && globalConfigFile.Config != nil {
		return globalConfigFile.Config, nil
	}

	if globalConfigFile == nil {
		globalConfigFile = &GlobalConfigFile{}
	}

	err := globalConfigFile.Load()
	if err != nil {
		return nil, err
	}

	return globalConfigFile.Config, nil
}

func SaveGlobalConfig() error {
	if globalConfigFile == nil {
		return errors.New("global config not loaded")
	}

	return globalConfigFile.Save()
}

type Secret struct {
	Name    string `json:"name" yaml:"name"`
	Use     string `json:"use" yaml:"use"`
	Key     string `json:"key" yaml:"key"`
	Vault   string `json:"vault" yaml:"vault"`
	Gen     bool   `json:"gen" yaml:"gen"`
	Special string `json:"special" yaml:"special"`
	Digits  bool   `json:"digits" yaml:"digits"`
	Lower   bool   `json:"lower" yaml:"lower"`
	Upper   bool   `json:"upper" yaml:"upper"`
	Size    int    `json:"size" yaml:"size"`
}

func (s *Secret) UnmarshalYAML(node *yaml.Node) error {
	s.Upper = true
	s.Lower = true
	s.Digits = true
	s.Size = 16
	s.Special = "#@!`~{}?.^&|=+_-"

	if node.Kind == yaml.ScalarNode {
		s.Name = node.Value
		return nil
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			value := node.Content[i+1]

			switch key.Value {
			case "name":
				s.Name = value.Value
			case "key":
				s.Key = value.Value
			case "use":
				s.Use = value.Value
			case "vault":
				s.Vault = value.Value
			case "gen":
				if strings.EqualFold(value.Value, "true") || value.Value == "1" {
					s.Gen = true
				} else if strings.EqualFold(value.Value, "false") || value.Value == "0" {
					s.Gen = false
				}
			case "special":
				s.Special = value.Value
			case "digits":
				if strings.EqualFold(value.Value, "true") || value.Value == "1" {
					s.Digits = true
				} else if strings.EqualFold(value.Value, "false") || value.Value == "0" {
					s.Digits = false
				}
			case "lower":
				if strings.EqualFold(value.Value, "true") || value.Value == "1" {
					s.Lower = true
				} else if strings.EqualFold(value.Value, "false") || value.Value == "0" {
					s.Lower = false
				}
			case "upper":
				if strings.EqualFold(value.Value, "true") || value.Value == "1" {
					s.Upper = true
				} else if strings.EqualFold(value.Value, "false") || value.Value == "0" {
					s.Upper = false
				}
			case "size":
				v, err := strconv.Atoi(value.Value)
				if err != nil {
					return err
				}

				s.Size = v
			}
		}
	}

	return nil
}

type Vault struct {
	Name string `json:"name" yaml:"name"`
	Uri  string `json:"uri" yaml:"uri"`
	Use  string `json:"use" yaml:"use"`
	With map[string]interface{}
}

type Secrets struct {
	order []string
	data  map[string]Secret
}

type Jolt9 struct {
	Name     string            `json:"name" yaml:"name"`
	Secrets  []Secret          `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Vaults   []Vault           `json:"vaults,omitempty" yaml:"vaults,omitempty"`
	Env      map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Dns      *Dns              `json:"dns,omitempty" yaml:"dns,omitempty"`
	Files    []string          `json:"files,omitempty" yaml:"files,omitempty"`
	Hooks    *Hooks            `json:"hooks,omitempty" yaml:"hooks,omitempty"`
	Inherits []string          `json:"inherits,omitempty" yaml:"inherits,omitempty"`
	Compose  *Compose          `json:"compose,omitempty" yaml:"compose,omitempty"`
	Ssh      *Ssh              `json:"ssh,omitempty" yaml:"ssh,omitempty"`
}

type Ssh struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Identity string `json:"identity" yaml:"identity"`
}

type Compose struct {
	Mode    string   `json:"mode" yaml:"mode"`
	Inline  string   `json:"inline,omitempty" yaml:"inline,omitempty"`
	Include []string `json:"include,omitempty" yaml:"include,omitempty"`
	Context string   `json:"context,omitempty" yaml:"context,omitempty"`
	Sudo    bool     `json:"sudo,omitempty" yaml:"sudo,omitempty"`
}

func (j *Jolt9) ResolveInheritence(dir string) (*Jolt9, error) {
	if len(j.Inherits) == 0 {
		return j, nil
	}

	var dest *Jolt9

	for _, inherit := range j.Inherits {
		if inherit == "" {
			continue
		}

		next, err := fs.Resolve(dir, inherit)
		if err != nil {
			return j, err
		}

		data, err := fs.ReadFile(next)
		if err != nil {
			return j, err
		}

		inheritFile := &Jolt9{}
		yaml.Unmarshal(data, inheritFile)

		if dest == nil {
			dest = inheritFile
			continue
		}

		dest.Merge(inheritFile)
	}

	dest.Merge(j)

	return dest, nil
}

func (j *Jolt9) PrependMergeVaults(vaults []Vault) {
	if len(vaults) == 0 {
		return
	}

	if j.Vaults == nil {
		j.Vaults = make([]Vault, 0)
	}

	if len(j.Vaults) == 0 {
		j.Vaults = vaults
		return
	}

	next := make([]Vault, 0)
	next = append(next, vaults...)

	for _, v := range j.Vaults {
		index := -1
		for i, v2 := range next {
			if v.Name == v2.Name {
				index = i
				break
			}
		}

		if index != -1 {
			next[index] = v
			continue
		}

		next = append(next, v)
	}

	j.Vaults = next
}

func (j *Jolt9) MergeVaults(vaults []Vault) {
	if len(vaults) == 0 {
		return
	}

	if j.Vaults == nil {
		j.Vaults = vaults
		return
	}

	for _, v := range vaults {
		index := -1
		for i, v2 := range j.Vaults {
			if v.Name == v2.Name {
				index = i
				break
			}
		}

		if index > -1 {
			j.Vaults[index] = v
			continue
		}

		j.Vaults = append(j.Vaults, v)
	}
}

func (j *Jolt9) PrependMergeSecrets(secrets []Secret) {
	if len(secrets) == 0 {
		return
	}

	if len(j.Secrets) == 0 {
		j.Secrets = secrets
		return
	}

	next := make([]Secret, 0)
	next = append(next, secrets...)

	for _, s := range j.Secrets {
		index := -1
		for i, s2 := range next {
			if s.Name == s2.Name {
				index = i
				break
			}
		}

		if index != -1 {
			next[index] = s
			continue
		}

		next = append(next, s)
	}

	j.Secrets = next
}

func (j *Jolt9) MergeSecrets(secrets []Secret) {
	if len(secrets) == 0 {
		return
	}

	if j.Secrets == nil {
		j.Secrets = secrets
		return
	}

	for _, s := range secrets {
		index := -1
		for i, s2 := range j.Secrets {
			if s.Name == s2.Name {
				index = i
				break
			}
		}

		if index > -1 {
			j.Secrets[index] = s
			continue
		}

		j.Secrets = append(j.Secrets, s)
	}
}

func (j *Jolt9) Merge(j2 *Jolt9) {
	if j2 == nil {
		return
	}

	if j2.Compose != nil {
		j.Compose = j2.Compose
	}

	if j2.Dns != nil {
		j.Dns = j2.Dns
	}

	j.MergeSecrets(j2.Secrets)
	j.MergeVaults(j2.Vaults)

	if j2.Dns != nil {
		if j.Dns == nil {
			j.Dns = j2.Dns
		} else if j2.Dns.Driver == "none" {
			j.Dns.Driver = "none"
			j.Dns.Env = map[string]string{}
			j.Dns.Zone = ""
		} else if j2.Dns.Use != "" {
			j.Dns.Use = j2.Dns.Use
			j.Dns.Driver = ""
			j.Dns.Env = map[string]string{}
			if j2.Dns.Zone != "" {
				j.Dns.Zone = j2.Dns.Zone
			}
		} else {
			if j2.Dns.Driver != "" {
				j.Dns.Driver = j2.Dns.Driver
			}

			if j2.Dns.Zone != "" {
				j.Dns.Zone = j2.Dns.Zone
			}

			if j2.Dns.Env != nil {
				if j.Dns.Env == nil {
					j.Dns.Env = map[string]string{}
				}

				for k, v := range j2.Dns.Env {
					j.Dns.Env[k] = v
				}
			}
		}
	}

	if len(j2.Files) > 0 {
		for _, f := range j2.Files {
			if slices.Contains(j.Files, f) {
				continue
			}

			j.Files = append(j.Files, f)
		}
	}

	if j2.Hooks != nil {
		if j2.Hooks.Before != nil && len(j2.Hooks.Before) > 0 {
			j.Hooks.Before = j2.Hooks.Before
		}

		if j2.Hooks.After != nil && len(j2.Hooks.After) > 0 {
			j.Hooks.After = j2.Hooks.After
		}
	}
}

type Dns struct {
	Driver string            `json:"driver" yaml:"driver"`
	Zone   string            `json:"zone" yaml:"zone"`
	Env    map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Use    string            `json:"use" yaml:"use"`
}

type Hooks struct {
	Before []Task `json:"before" yaml:"before"`
	After  []Task `json:"after" yaml:"after"`
}

type Task struct {
	Name    string            `json:"name" yaml:"name"`
	Run     string            `json:"run" yaml:"run"`
	Timeout string            `json:"timeout" yaml:"timeout"`
	Env     map[string]string `json:"env" yaml:"env"`
	Use     string            `json:"use" yaml:"use"`
}

type Workspace struct {
	Name      string             `json:"name" yaml:"name"`
	Vaults    []Vault            `json:"vaults,omitempty" yaml:"vaults,omitempty"`
	Secrets   []Secret           `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	Dns       map[string]Dns     `json:"dns,omitempty" yaml:"dns,omitempty"`
	Projects  map[string]Project `json:"projects" yaml:"projects"`
	Discovery *Discovery         `json:"discovery,omitempty" yaml:"discovery,omitempty"`
}

type Project struct {
	Targets map[string]string `json:"targets" yaml:"targets"`
}

func (p *Project) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		p.Targets = map[string]string{
			"default": node.Value,
		}

		return nil
	}

	if node.Kind == yaml.MappingNode {
		p.Targets = make(map[string]string)
		for i := 0; i < len(node.Content); i += 2 {
			key := node.Content[i]
			value := node.Content[i+1]

			p.Targets[key.Value] = value.Value
		}
	}

	return nil
}

type Discovery struct {
	Include []string `json:"include" yaml:"include"`
	Exclude []string `json:"exclude" yaml:"exclude"`
}

type WorkspaceFile struct {
	File   string
	Config *Workspace
}

func (cfg *WorkspaceFile) Load() error {
	cfg.Config = &Workspace{}

	data, err := fs.ReadFile(cfg.File)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, cfg.Config)
}

func (cfg *WorkspaceFile) Save() error {
	data, err := yaml.Marshal(cfg.Config)
	if err != nil {
		return err
	}

	return fs.WriteFile(cfg.File, data, 0644)
}

type GlobalConfig struct {
	Hosts      map[string]Host   `json:"hosts" yaml:"hosts"`
	Workspaces map[string]string `json:"scopes" yaml:"scopes"`
	Paths      *GlobalPaths      `json:"paths" yaml:"paths"`
}

type GlobalPaths struct {
	Cache string `json:"cache" yaml:"cache"`
}

type GlobalConfigFile struct {
	File   string
	Config *GlobalConfig
}

func (cfg *GlobalConfigFile) Load() error {
	cfgDir, err := paths.ConfigDir()
	if err != nil {
		return err
	}

	if cfg.Config == nil {
		cfg.Config = &GlobalConfig{}
	}

	if !fs.Exists(cfgDir) {
		return nil
	}

	try := []string{"config.yaml", "config.yml", "config.json"}
	for _, t := range try {
		file := filepath.Join(cfgDir, t)
		if fs.Exists(file) {
			cfg.File = file
			data, err := fs.ReadFile(file)
			if err != nil {
				return err
			}

			ext := filepath.Ext(file)
			switch ext {
			case ".json":
				err = json.Unmarshal(data, cfg.Config)
			case ".yaml", ".yml":
				err = yaml.Unmarshal(data, cfg.Config)
			}

			if err != nil {
				return err
			}

			return nil
		}
	}

	if cfg.Config.Paths == nil {
		cfg.Config.Paths = &GlobalPaths{}
		if cfg.Config.Paths.Cache == "" {
			cacheDir, err := paths.CacheDir()
			if err != nil {
				return err
			}

			cfg.Config.Paths.Cache = cacheDir
		}
	}

	if cfg.Config.Workspaces == nil {
		cfg.Config.Workspaces = make(map[string]string)
	}

	_, ok := cfg.Config.Workspaces["@default"]
	if !ok {
		dir, err := paths.ConfigDir()
		if err != nil {
			return err
		}

		// global default workspace.
		wsf := filepath.Join(dir, consts.WorkspaceFileName)

		cfg.Config.Workspaces["default"] = wsf

		ws := Workspace{
			Name: "default",
		}

		ws.Discovery = &Discovery{
			Include: []string{},
			Exclude: []string{},
		}

		ws.Projects = make(map[string]Project)

		wsFile := &WorkspaceFile{
			File:   wsf,
			Config: &ws,
		}

		err = wsFile.Save()
		if err != nil {
			return err
		}
	}

	return nil
}

func (cfg *GlobalConfigFile) Save() error {
	if cfg.File == "" {
		cfgDir, err := paths.ConfigDir()
		if err != nil {
			return err
		}

		cfg.Config = &GlobalConfig{}

		err = fs.EnsureDir(cfgDir, 0755)
		if err != nil {
			return err
		}

		cfg.File = filepath.Join(cfgDir, "config.yaml")
	}

	ext := filepath.Ext(cfg.File)
	switch ext {
	case ".json":
		data, err := json.MarshalIndent(cfg.Config, "", "    ")
		if err != nil {
			return err
		}

		return fs.WriteFile(cfg.File, data, 0644)
	case ".yaml", ".yml":
		data, err := yaml.Marshal(cfg.Config)
		if err != nil {
			return err
		}

		return fs.WriteFile(cfg.File, data, 0644)
	}

	return errors.New("unsupported file extension")
}

type Host struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Identity string `json:"identity" yaml:"identity"`
}
