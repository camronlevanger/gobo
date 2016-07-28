package commands

import (
	"errors"

	"github.com/camronlevanger/gobo/utils"
)

// IDeleteCommand is the interface to implement when deleting any gobo environments.
type IDeleteCommand interface {
	Run(name string)
}

// DeleteCommand is the struct for this implementation of IDeleteCommand.
type DeleteCommand struct {
	logger        utils.ILogger
	moveService   utils.IMoveService
	configService utils.IConfigService
	gobopath      string
	gopath        string
}

// GetDeleteCommand returns a pointer to an implementation of IDeleteCommand.
func GetDeleteCommand(
	logger utils.ILogger,
	moveService utils.IMoveService,
	configService utils.IConfigService,
	gobopath string,
	gopath string,
) *DeleteCommand {
	delete := DeleteCommand{
		logger,
		moveService,
		configService,
		gobopath,
		gopath,
	}

	return &delete
}

// Run deletes the virtual environment 'name'.
func (delete *DeleteCommand) Run(name string) error {

	env := delete.configService.ReadEnvironment(delete.gopath + "gobo.toml")

	if env.Name == name {
		return errors.New("You may not delete the active environment.")
	}

	delete.logger.Info("Removing virtual environment " + name + " at " + delete.gobopath + name)
	err := delete.moveService.RemoveDirectory(delete.gobopath + name)

	return err
}
