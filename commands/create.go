package commands

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/camronlevanger/gobo/models"
	"github.com/camronlevanger/gobo/utils"
)

// ICreateCommand defines the interface for creating new environments.
type ICreateCommand interface {
	Run(name string) error
}

// CreateCommand is the struct for this instance of ICreateCommand interface.
type CreateCommand struct {
	logger         utils.ILogger
	configService  utils.IConfigService
	copyService    utils.ICopyService
	moveService    utils.IMoveService
	packageService utils.IPackageService
	populate       bool
	gopath         string
	gobopath       string
	host           models.Host
	initial        bool
}

// GetCreateCommand returns an implementation of ICreateCommand.
func GetCreateCommand(
	logger utils.ILogger,
	configService utils.IConfigService,
	copyService utils.ICopyService,
	moveService utils.IMoveService,
	packageService utils.IPackageService,
	populate bool,
	gopath string,
	gobopath string,
	host models.Host,
	initial bool,
) *CreateCommand {
	var create = CreateCommand{
		logger,
		configService,
		copyService,
		moveService,
		packageService,
		populate,
		gopath,
		gobopath,
		host,
		initial,
	}

	return &create
}

// Run creates a new virtual environment with the name provided.
func (create *CreateCommand) Run(name string) error {

	var current = "initial"

	_, existsErr := os.Stat(create.gopath + "gobo.toml")
	if existsErr == nil {
		create.logger.Info("A current environment seems to exist...")

		env := create.configService.ReadEnvironment(create.gopath + "gobo.toml")
		current = env.Name

		// make sure we aren't creating the environment we are in
		if env.Name == name {
			return errors.New(name + " is already the currently active environment.")
		}

		// make sure we aren't creating an existing environment
		envs, _ := ioutil.ReadDir(create.gobopath)
		count := 0
		for _, e := range envs {
			if e.IsDir() {
				count++
				if e.Name() == name {
					return errors.New(name + " is already a named environment.")
				}
			}
		}

		// if there is a current env, try to save it first.
		if !create.initial {

			create.logger.Info("Running save on current environment first...")

			save := GetSaveCommand(
				create.logger,
				create.configService,
				create.packageService,
				create.copyService,
				create.host,
				create.gopath,
				create.gobopath,
			)

			err := save.Run(false)
			if err != nil {
				create.logger.Error("Error running gobo save command: " + err.Error())
			}
		}
	}

	var installedPackages []models.Package

	if !create.populate {
		for _, dir := range models.GOPATHDIRECTORIES {
			create.logger.Info(fmt.Sprintf("Removing active %s directory at %s\n", dir, create.gopath))
			err := create.moveService.Move(create.gopath+dir, create.gobopath+current)
			if err != nil {
				create.logger.Error("Error moving GOPATH directory: " + err.Error())
			}
		}
	} else {
		create.logger.Info("Populating this new environment, so looking up installed packages and their bookmarks...")

		installedPackages = create.packageService.GetInstalledPackages()
	}

	if !create.initial {
		err := create.moveService.Move(create.gopath+"gobo.toml", create.gobopath+current)
		if err != nil {
			create.logger.Error("Error moving current toml env file: " + err.Error())
		}

		err = create.moveService.Move(create.gopath+"packages.toml", create.gobopath+current)
		if err != nil {
			create.logger.Error("Error moving current toml packages file: " + err.Error())
		}
	}

	environment := models.Environment{}
	environment.Name = name
	environment.Host = create.host
	environment.DateCreated = time.Now()
	environment.DateModified = time.Now()

	deps := models.Dependencies{}
	deps.Package = installedPackages

	envpath := create.gopath + "gobo.toml"
	pakpath := create.gopath + "packages.toml"

	create.logger.Info("Writing environment file (gobo.toml).")
	err := create.configService.WriteEnvironment(envpath, environment)

	create.logger.Info("Writing packages file (packages.toml).")
	err = create.configService.WritePackages(pakpath, deps)

	if err != nil {
		return err
	}

	err = os.Mkdir(create.gobopath+name, models.FILEMODE)
	if err != nil {
		return err

	}

	_, existsErr = os.Stat(create.gopath + "bin")
	if existsErr != nil {
		os.MkdirAll(create.gopath+"bin", models.FILEMODE)
	}

	_, existsErr = os.Stat(create.gopath + "pkg")
	if existsErr != nil {
		os.MkdirAll(create.gopath+"pkg", models.FILEMODE)
	}

	_, existsErr = os.Stat(create.gopath + "src")
	if existsErr != nil {
		os.MkdirAll(create.gopath+"src", models.FILEMODE)
	}

	if !create.populate {
		err = create.copyService.CopyFile(create.gobopath+"gobo", create.gopath+"bin")
	}
	return err
}
