package config

import (
	"errors"
	"fmt"
	"os"
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
	yamlConf, err := os.ReadFile(abspath)
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

func ParseDBConnConfigEnv() (*DBConnConfig, error) {
	var dbconf DBConnConfig
	var ok bool
	dbconf.DB_USER, ok = os.LookupEnv("DB_USER")
	if !ok {
		return nil, errors.New("db user not specified")
	}
	dbconf.DB_PASSWORD, ok = os.LookupEnv("DB_PASSWORD")
	if !ok {
		return nil, errors.New("db password not specified")
	}
	dbconf.DB_NAME, ok = os.LookupEnv("DB_NAME")
	if !ok {
		return nil, errors.New("db name not specified")
	}
	dbconf.DB_HOST, ok = os.LookupEnv("DB_HOST")
	if !ok {
		return nil, errors.New("db host not specified")
	}
	dbconf.DB_PORT, ok = os.LookupEnv("DB_PORT")
	if !ok {
		return nil, errors.New("db port not specified")
	}
	return &dbconf, nil
}
