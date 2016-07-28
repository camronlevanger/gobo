package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/camronlevanger/gobo/models"
)

// IConfigService is the interface to implement for reading and writing environments.
type IConfigService interface {
	WriteEnvironment(path string, env models.Environment) error
	ReadEnvironment(path string) models.Environment
	WritePackages(path string, paks models.Dependencies) error
	ReadPackages(path string) models.Dependencies
}

// ConfigService is the struct for this implementation of IConfigService.
type ConfigService struct {
	logger ILogger
}

// GetConfigService returns a pointer to an implementation of IConfigService.
func GetConfigService(logger ILogger) IConfigService {
	var configService = ConfigService{
		logger,
	}

	return &configService
}

// WriteEnvironment writes out an environment struct to the given file location.
func (configService *ConfigService) WriteEnvironment(path string, env models.Environment) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(env); err != nil {
		configService.logger.Fatal(err.Error())
	}
	configService.logger.Info(fmt.Sprintf("Writing environment to %s:\n", path))
	//configService.logger.Info(fmt.Sprint(buf.String()))

	err := ioutil.WriteFile(path, []byte(buf.String()), 0755)

	return err
}

// ReadEnvironment loads the toml file at the provided path into an Environment struct and returns it for use.
func (configService *ConfigService) ReadEnvironment(path string) models.Environment {

	var env models.Environment

	if _, err := toml.DecodeFile(path, &env); err != nil {
		configService.logger.Fatal(fmt.Sprintf("Unable to read environment file at %s because: %s", path, err.Error()))
	}

	return env
}

// ReadPackages loads the toml file at the provided path into a Dependencies struct and returns it for use.
func (configService *ConfigService) ReadPackages(path string) models.Dependencies {

	var paks models.Dependencies

	if _, err := toml.DecodeFile(path, &paks); err != nil {
		configService.logger.Fatal(fmt.Sprintf("Unable to read packages file at %s because: %s", path, err.Error()))
	}

	return paks
}

// WritePackages writes out a Dependencies struct to the given file location.
func (configService *ConfigService) WritePackages(path string, paks models.Dependencies) error {

	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(paks); err != nil {
		configService.logger.Fatal(err.Error())
	}
	configService.logger.Info(fmt.Sprintf("Writing packages to %s:\n", path))
	//configService.logger.Info(fmt.Sprint(buf.String()))

	err := ioutil.WriteFile(path, []byte(buf.String()), 0755)

	return err
}
