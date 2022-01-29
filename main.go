package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Port int
}

type handler struct {
	service BusinessService
}

func (h handler) getServices(c *gin.Context) {

	services, err := h.service.GetServices()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, services)
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

	repo := NewTestingRepo()

	service := NewBusinessService(repo)

	h := handler{
		service: service,
	}

	r := gin.Default()

	r.GET("/services", h.getServices)

	portString := strconv.FormatInt(int64(config.Port), 10)

	os.Setenv("PORT", portString)

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
