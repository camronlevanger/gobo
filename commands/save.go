package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/camronlevanger/gobo/models"
	"github.com/camronlevanger/gobo/utils"
)

// ISaveCommand is the interface to implement for saving changes to the environment.
type ISaveCommand interface {
	Run(silent bool) error
}

// SaveCommand is the struct for this implementation of ISaveCommand.
type SaveCommand struct {
	logger         utils.ILogger
	configService  utils.IConfigService
	packageService utils.IPackageService
	copyService    utils.ICopyService
	host           models.Host
	gopath         string
	gobopath       string
}

// GetSaveCommand returns a pointer to an implementation of the ISaveCommand interface.
func GetSaveCommand(
	logger utils.ILogger,
	configService utils.IConfigService,
	packageService utils.IPackageService,
	copyService utils.ICopyService,
	host models.Host,
	gopath string,
	gobopath string,
) *SaveCommand {
	save := SaveCommand{
		logger,
		configService,
		packageService,
		copyService,
		host,
		gopath,
		gobopath,
	}

	return &save
}

// Run diffs the environment file to the filesystem, and confirms writing the detected changes if found.
func (save *SaveCommand) Run(silent bool) error {

	envFile := save.gopath + "gobo.toml"
	pakFile := save.gopath + "packages.toml"

	env := save.configService.ReadEnvironment(envFile)
	pak := save.configService.ReadPackages(pakFile)

	changed, updatedPackages := save.packageService.DiffAndUpdatePackages(pak.Package)
	env.Host = save.host
	pak.Package = updatedPackages

	if changed {
		if !silent {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("The current environment has uncommited changes, update now? (y): ")
			answer, _ := reader.ReadString('\n')
			if answer == "n\n" || answer == "N\n" {
				save.logger.Info("Not saving environment updates.")
				return nil
			}

			save.logger.Info("Commiting updates to " + env.Name)
			save.configService.WriteEnvironment(save.gopath, env)
			save.configService.WritePackages(save.gopath, pak)
		} else {
			save.configService.WriteEnvironment(save.gopath, env)
			save.configService.WritePackages(save.gopath, pak)
		}
	} else {
		save.logger.Info("No environment changes detected.")
	}

	return nil
}
