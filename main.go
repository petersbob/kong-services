package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port int
}

func main() {

	config, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = validateConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	portString := strconv.FormatInt(int64(config.Port), 10)

	os.Setenv("PORT", portString)

	r := gin.Default()

	r.GET("/greeting", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"hello": "world",
		})
	})

	r.Run()

}

func validateConfig(config Config) error {
	if config.Port == 0 {
		return errors.New("missing port config value")
	}

	return nil
}

func getConfig() (Config, error) {
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return Config{}, err
	}

	config := Config{}

	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}
