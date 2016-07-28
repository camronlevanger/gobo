package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/camronlevanger/gobo/models"
	"github.com/camronlevanger/gobo/utils"
)

// IRestoreCommand is the interface to implement for restoring the system to its pre-gobo state.
type IRestoreCommand interface {
	Run() error
}

// RestoreCommand is the struct for this implementation of IRestoreCommand.
type RestoreCommand struct {
	logger      utils.ILogger
	copyService utils.ICopyService
	gopath      string
	gobopath    string
}

// GetRestoreCommand returns a pointer to an implementation of IRestoreCommand.
func GetRestoreCommand(
	logger utils.ILogger,
	copyService utils.ICopyService,
	gopath string,
	gobopath string,
) *RestoreCommand {
	var restore = RestoreCommand{
		logger,
		copyService,
		gopath,
		gobopath,
	}

	return &restore
}

// Run deletes all traces of gobo and restores the original gopath from ~/.gobo/initial.
func (restore *RestoreCommand) Run() error {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("The restore command deletes all virtual environments and gobo files, then restores your GOPATH to the state it was in before you ran gobo for the first time. Are you sure this is what you want to do? (yes/no): ")
	answer, _ := reader.ReadString('\n')
	if answer == "yes\n" || answer == "YES\n" {
		for _, dir := range models.GOPATHDIRECTORIES {
			restore.logger.Error(fmt.Sprintf("Removing active %s directory at %s\n", dir, restore.gopath))
			err := os.RemoveAll(restore.gopath + dir)
			if err != nil {
				restore.logger.Error("RESTORE - Error cleaning up current GOPATH: " + err.Error())
			}
		}

		for _, dir := range models.GOPATHDIRECTORIES {
			restore.logger.Error(fmt.Sprintf("Restoring %s directory from %s\n", dir, restore.gobopath+dir))
			err := restore.copyService.CopyDir(restore.gobopath+"initial"+dir, restore.gopath)
			if err != nil {
				restore.logger.Error("RESTORE - Error copying initial backup directory to GOPATH: " + err.Error())
			}

		}

		os.RemoveAll(restore.gobopath)
	}

	return nil
}
