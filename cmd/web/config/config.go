package config

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
)

type DBConnConfig struct {
	Host     string `env:"HOST,required"`
	Port     string `env:"PORT,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Name     string `env:"NAME,required"`
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

func ParseDBConnConfigEnv(ctx context.Context, prefix string) (*DBConnConfig, error) {
	var dbconf DBConnConfig
	l := envconfig.PrefixLookuper(prefix, envconfig.OsLookuper())
	err := envconfig.ProcessWith(ctx, &dbconf, l)
	if err != nil {
		return nil, fmt.Errorf("error processing: %w", err)
	}
	return &dbconf, err
}
