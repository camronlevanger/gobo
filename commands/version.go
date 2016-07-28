package commands

import (
	"fmt"

	"github.com/camronlevanger/gobo/models"
)

// IVersionCommand is the interface to implememt for getting the gobo version number.
type IVersionCommand interface {
}

// VersionCommand is the struct for this implementation of IVersionCommand.
type VersionCommand struct {
}

// GetVersionCommand returns a pointer to an implementation of IVersionCommand.
func GetVersionCommand() *VersionCommand {
	version := VersionCommand{}

	return &version
}

// Run prints the version constant to the console.
func (version *VersionCommand) Run() {
	fmt.Println("gobo version " + models.GOBOVERSION)
}
