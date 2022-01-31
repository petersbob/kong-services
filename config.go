package main

import (
	"errors"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type config struct {
	Port           int    `yaml:"port"`
	DatabaseURL    string `yaml:"databaseurl"`
	MigrationsPath string `yaml:"migrationspath"`
}

func validateConfig(config config) error {
	if config.Port == 0 {
		return errors.New("missing port config value")
	}

	if strings.TrimSpace(config.DatabaseURL) == "" {
		return errors.New("missing database url config value")
	}

	if strings.TrimSpace(config.MigrationsPath) == "" {
		return errors.New("missing migrations path config value")
	}

	return nil
}

func getConfig() (config, error) {
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return config{}, err
	}

	currentConfig := config{}

	err = yaml.Unmarshal(configFile, &currentConfig)
	if err != nil {
		return config{}, err
	}

	return currentConfig, nil
}
