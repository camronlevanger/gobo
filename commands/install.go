package commands

import (
	"github.com/camronlevanger/gobo/models"
	"github.com/camronlevanger/gobo/utils"
)

// IInstallCommand is the interface to install go packages at a specific hash.
type IInstallCommand interface {
	Run(file string) error
}

// InstallCommand is the struct for this implementation of IInstallCommand.
type InstallCommand struct {
	logger         utils.ILogger
	configService  utils.IConfigService
	packageService utils.IPackageService
	copyService    utils.ICopyService
	host           models.Host
	gopath         string
	gobopath       string
}

// GetInstallCommand returns a pointer to an implmentation of IInstallCommand.
func GetInstallCommand(
	logger utils.ILogger,
	configService utils.IConfigService,
	packageService utils.IPackageService,
	copyService utils.ICopyService,
	host models.Host,
	gopath string,
	gobopath string,
) *InstallCommand {
	install := InstallCommand{
		logger,
		configService,
		packageService,
		copyService,
		host,
		gopath,
		gobopath,
	}

	return &install
}

// Run loops through all packages in the install file and then runs get, checkout, install utils on them.
func (install *InstallCommand) Run(file string) error {
	paks := install.configService.ReadPackages(file)

	install.logger.Info("Installing packages from environment file: " + file)

	packages := paks.Package
	var err error

	for i := 0; i < len(packages); i++ {
		err = install.packageService.Get(packages[i].Path, packages[i].Revision)
		if err != nil {
			install.logger.Error("Error installing " + packages[i].Path + " because: " + err.Error())
		}
	}

	install.logger.Info("Saving the environment.")
	save := GetSaveCommand(
		install.logger,
		install.configService,
		install.packageService,
		install.copyService,
		install.host,
		install.gopath,
		install.gobopath,
	)

	err = save.Run(true)

	return err
}
