package main

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	Port int
}

func validateConfig(config config) error {
	if config.Port == 0 {
		return errors.New("missing port config value")
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
