package utils

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

// ICopyService is the interface to implement for cross platform copy commands.
type ICopyService interface {
	CopyDir(source string, destination string) error
	CopyFile(source string, destination string) error
}

// CopyService is the struct for this implementation of the ICopyService interface.
type CopyService struct {
}

// GetCopyService returns a pointer to an implementation of the ICopyService interface.
func GetCopyService() *CopyService {
	var copy = CopyService{}

	return &copy
}

// CopyDir uses system copy commands to recursively copy directory structures.
func (copy *CopyService) CopyDir(source string, destination string) error {
	var cmd string
	var err error
	if runtime.GOOS == "windows" {
		fmt.Println("Hello from Windows")
		cmd, err := exec.LookPath("xcopy")
		if err != nil {
			return errors.New(fmt.Sprintf("Gobo for Windows requires xcopy: %s\n", err.Error()))
		}

		cmdArgs := []string{"/E", source, destination}
		if _, err = exec.Command(cmd, cmdArgs...).Output(); err != nil {
			return errors.New(fmt.Sprintf(
				"There was an error making a copy of the current GOPATH for backup: %s\n",
				err,
			),
			)
		}
	} else {
		cmd = "cp"
		cmdArgs := []string{"-r", source, destination}

		cmd := exec.Command(cmd, cmdArgs...)

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()

		if err != nil {
			return errors.New(
				"There was an error making a copy of the current GOPATH for backup: " +
					err.Error() + ": " +
					stderr.String(),
			)
		}
	}

	return err

}

// CopyFile uses system copy commands to recursively copy directory structurescopy files.
func (copy *CopyService) CopyFile(source string, destination string) error {
	var cmd string
	var err error
	if runtime.GOOS == "windows" {
		fmt.Println("Hello from Windows")
		cmd, err := exec.LookPath("xcopy")
		if err != nil {
			return errors.New(fmt.Sprintf("Gobo for Windows requires xcopy: %s\n", err.Error()))
		}

		cmdArgs := []string{source, destination}
		if _, err = exec.Command(cmd, cmdArgs...).Output(); err != nil {
			return errors.New(fmt.Sprintf(
				"There was an error making a copy of the current GOPATH for backup: %s\n",
				err,
			),
			)
		}
	} else {
		cmd = "cp"
		cmdArgs := []string{source, destination}
		if _, err = exec.Command(cmd, cmdArgs...).Output(); err != nil {
			return errors.New(fmt.Sprintf(
				"There was an error making a copy of the current GOPATH for backup: %s\n",
				err,
			))
		}
	}

	return err

}
