package internal

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name   string `yaml:"name"`
	Dir    string `yaml:"dir"`
	CmdStr string `yaml:"cmdStr"`
	ProxyPort  string `yaml:"proxyPort"`
}

type configFile struct {
	Servers []Config `yaml:"servers"`
}

func LoadConfig(path string) ([]Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg configFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return cfg.Servers, nil
}
