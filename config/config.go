package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Username         string `yaml:"user"`
		Password         string `yaml:"pass"`
		ConnectionString string `yaml:"conn"`
		DatabaseIP       string `yaml:"dbip"`
	} `yaml:"database"`
}

func GetConfig() (*Config, error) {
	f, err := os.Open("config/config.yml")
	if err != nil {
		fmt.Print(err)
		return nil, fmt.Errorf("error while reading file")
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, fmt.Errorf("error while decoding file")
	}
	return &cfg, nil
}
