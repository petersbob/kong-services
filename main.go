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

func (h handler) getService(c *gin.Context) {
	typeCodeString := c.Param("type_code")

	if typeCodeString == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid type code"))
		return
	}

	typeCode, err := strconv.ParseUint(typeCodeString, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid type code"))
		return
	}

	service, err := h.service.GetService(ServiceTypeCode(typeCode))
	if err != nil {
		if err == errorNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, service)
}

func (h handler) getServiceVersions(c *gin.Context) {
	typeCodeString := c.Param("type_code")

	if typeCodeString == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid type code"))
		return
	}

	typeCode, err := strconv.ParseUint(typeCodeString, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid type code"))
		return
	}

	serviceVersions, err := h.service.GetServiceVersions(ServiceTypeCode(typeCode))
	if err != nil {
		if err == errorNotFound {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, serviceVersions)
}

func (h handler) getServiceVersion(c *gin.Context) {
	typeCodeString := c.Param("type_code")

	if typeCodeString == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid type code"))
		return
	}

	typeCode, err := strconv.ParseUint(typeCodeString, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid type code"))
		return
	}

	versionNumberString := c.Param("version_number")
	if versionNumberString == "" {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid version number"))
		return
	}

	versionNumber, err := strconv.ParseUint(versionNumberString, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, errors.New("invalid version number"))
	}

	serviceVersions, err := h.service.GetServiceVersion(ServiceTypeCode(typeCode), uint(versionNumber))
	if err != nil {
		if err == errorNotFound {
			c.AbortWithError(http.StatusNotFound, errors.New("service version not found"))
			return
		}

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, serviceVersions)
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
	r.GET("/services/:type_code", h.getService)
	r.GET("/services/:type_code/versions", h.getServiceVersions)
	r.GET("/services/:type_code/versions/:version_number", h.getServiceVersion)

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
