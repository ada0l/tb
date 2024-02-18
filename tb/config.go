package tb

import (
	"gopkg.in/yaml.v3"
	"os"
)

type TbConfig struct {
	Host     string   `yaml:"host"`
	Backends []string `yaml:"backends"`
}

func LoadConfig(configBytes []byte) (config *TbConfig, err error) {
	config = &TbConfig{}
	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		return
	}
	return
}

func LoadConfigFromFile(configPath string) (config *TbConfig, err error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return
	}
	config, err = LoadConfig(content)
	return
}
