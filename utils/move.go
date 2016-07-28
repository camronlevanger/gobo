package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

// IMoveService is the interface tha exposes execs to system mv cmd.
type IMoveService interface {
	Move(source string, dest string) error
	RemoveDirectory(source string) error
}

// MoveService is the struct for this instance of IMoveService.
type MoveService struct {
}

// GetMoveService returns a pointer to an implementation of IMoveService.
func GetMoveService() *MoveService {
	move := MoveService{}

	return &move
}

// Move moves the source path to the dest path.
func (move *MoveService) Move(source string, dest string) error {
	var app string
	var err error
	if runtime.GOOS == "windows" {
		app = "move"
		cmdArgs := []string{"-f", source, dest}
		if _, err = exec.Command(app, cmdArgs...).Output(); err != nil {
			return errors.New(fmt.Sprintf("There was an error moving the current environment: %s\n",
				err,
			),
			)
		}
	} else {
		cmdArgs := []string{"-f", source, dest}
		app = "mv"

		cmd := exec.Command(app, cmdArgs...)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()

		if err != nil {
			fmt.Println("Error running mv: " + err.Error() + ": " + stderr.String())
			return err
		}

		return err
	}

	return err
}

// RemoveDirectory moves a directory to nowhere.
func (move *MoveService) RemoveDirectory(source string) error {
	var app string
	var err error
	if runtime.GOOS == "windows" {
		app = "rmdir"
		cmdArgs := []string{"/s", source}
		if _, err = exec.Command(app, cmdArgs...).Output(); err != nil {
			return errors.New(fmt.Sprintf("There was an error moving the current environment: %s\n",
				err,
			),
			)
		}
	} else {
		cmdArgs := []string{"-rf", source}
		app = "rm"

		cmd := exec.Command(app, cmdArgs...)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()

		if err != nil {
			fmt.Println("Error running rm -rf: " + err.Error() + ": " + stderr.String())
			return err
		}

		return err
	}

	return err
}
