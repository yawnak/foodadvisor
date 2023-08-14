package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
)

type DBConnConfig struct {
	Host     string `config:"dbhost"`
	Port     string `config:"dbport"`
	User     string `config:"dbuser"`
	Name     string `config:"dbname"`
	Password string
}

func BindConfig(conf any, paths ...string) error {
	defer config.ClearAll()
	if reflect.ValueOf(conf).Type().Kind() != reflect.Ptr {
		return errors.New("conf is not a pointer")
	}
	config.WithOptions(config.ParseEnv, func(opt *config.Options) {
		opt.DecoderConfig.TagName = "config"
	})
	// add driver for support yaml content
	config.AddDriver(yaml.Driver)

	err := config.LoadFiles(paths...)
	if err != nil {
		return fmt.Errorf("error loading config files: %w", err)
	}

	err = config.Decode(conf)
	return err
}
