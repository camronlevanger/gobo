package commands

import (
	"errors"
	"fmt"
	"os"

	"github.com/camronlevanger/gobo/models"
	"github.com/camronlevanger/gobo/utils"
)

type IActivateCommand interface {
}

type ActivateCommand struct {
	logger         utils.ILogger
	configService  utils.IConfigService
	packageService utils.IPackageService
	copyService    utils.ICopyService
	moveService    utils.IMoveService
	host           models.Host
	gopath         string
	gobopath       string
}

func GetActivateCommand(
	logger utils.ILogger,
	configService utils.IConfigService,
	packageService utils.IPackageService,
	copyService utils.ICopyService,
	moveService utils.IMoveService,
	host models.Host,
	gopath string,
	gobopath string,
) *ActivateCommand {
	activate := ActivateCommand{
		logger,
		configService,
		packageService,
		copyService,
		moveService,
		host,
		gopath,
		gobopath,
	}

	return &activate
}

func (activate *ActivateCommand) Run(name string) error {

	env := activate.configService.ReadEnvironment(activate.gopath + "gobo.toml")

	if env.Name == name {
		return errors.New(name + " is already the currently active environment.")
	}

	_, existsErr := os.Stat(activate.gopath + "gobo.toml")
	if existsErr == nil {

		activate.logger.Info("Running save on current environment first...")

		save := GetSaveCommand(
			activate.logger,
			activate.configService,
			activate.packageService,
			activate.copyService,
			activate.host,
			activate.gopath,
			activate.gobopath,
		)

		err := save.Run(false)
		if err != nil {
			activate.logger.Error("Error running gobo save command: " + err.Error())
		}
	}

	for _, dir := range models.GOPATHDIRECTORIES {
		activate.logger.Info(fmt.Sprintf("Removing active %s directory at %s\n", dir, activate.gopath))
		activate.moveService.Move(activate.gopath+dir, activate.gobopath+env.Name)
	}

	activate.logger.Info("Removing gobo.toml...")
	activate.moveService.Move(activate.gopath+"gobo.toml", activate.gobopath+env.Name)

	activate.logger.Info("Removing packages.toml...")
	activate.moveService.Move(activate.gopath+"packages.toml", activate.gobopath+env.Name)

	activate.logger.Info("Activating " + name)
	for _, dir := range models.GOPATHDIRECTORIES {
		activate.logger.Info(fmt.Sprintf("Moving %s to %s\n", activate.gobopath+name+"/"+dir, activate.gopath))
		activate.moveService.Move(activate.gobopath+name+"/"+dir, activate.gopath)
	}

	activate.logger.Info("Restoring gobo.toml...")
	activate.moveService.Move(activate.gobopath+name+"/gobo.toml", activate.gopath)

	activate.logger.Info("Restoring packages.toml...")
	activate.moveService.Move(activate.gobopath+name+"/packages.toml", activate.gopath)

	return nil
}
