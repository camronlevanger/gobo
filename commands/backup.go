package commands

import (
	"fmt"
	"os"

	"github.com/camronlevanger/gobo/models"
	"github.com/camronlevanger/gobo/utils"
)

// IBackupCommand is the interface to implement for creating backups of a GOPATH.
type IBackupCommand interface {
	Run() (bool, error)
}

// BackupCommand is the struct for this implementation of the IBackupCommand interface.
type BackupCommand struct {
	logger      utils.ILogger
	copyService utils.ICopyService
	home        string
	gopath      string
	gobo        string
	restore     string
	goboMaster  string
}

// GetBackupCommand returns a pointer to an implementation of the IBackupCommand interface.
func GetBackupCommand(
	logger utils.ILogger,
	copyService utils.ICopyService,
	home string,
	gopath string,
	gobo string,
	restore string,
	goboMaster string,
) *BackupCommand {
	var backup = BackupCommand{
		logger,
		copyService,
		home,
		gopath,
		gobo,
		restore,
		goboMaster,
	}

	return &backup
}

// Run copies the current GOPATH into a gobo backup directory.
func (backup *BackupCommand) Run() (bool, error) {
	var err error
	var initial bool

	if _, err = os.Stat(backup.restore); err == nil {
		backup.logger.Info("Gobo initial backup already exists.")
	} else {
		backup.logger.Info(fmt.Sprintf("Creating Gobo initial backup at %s...\n", backup.restore))

		os.MkdirAll(backup.restore, models.FILEMODE)

		err = backup.copyService.CopyDir(backup.gopath, backup.restore)

		err = backup.copyService.CopyFile(backup.gopath+"bin/gobo", backup.gobo)

		initial = true

	}
	return initial, err
}
