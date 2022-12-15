package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type DBConnConfig struct {
	DB_HOST     string `yaml:"dbhost"`
	DB_PORT     string `yaml:"dbport"`
	DB_USER     string `yaml:"dbuser"`
	DB_PASSWORD string `yaml:"dbpassword"`
	DB_NAME     string `yaml:"dbname"`
}

func ParseDBConnConfig(path string) (*DBConnConfig, error) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("bad path: %w", err)
	}
	yamlConf, err := ioutil.ReadFile(abspath)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	var dbconf DBConnConfig
	err = yaml.Unmarshal(yamlConf, &dbconf)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}
	return &dbconf, nil
}
