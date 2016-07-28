package models

import (
	"os/user"
	"time"
)

// Environment is a struct representing a gobo TOML configuration file.
type Environment struct {
	Name         string    `toml:"name"`
	DateCreated  time.Time `toml:"created"`
	DateModified time.Time `toml:"modified"`
	Host         Host      `toml:"host"`
}

// Host is a struct describing User and System information that acted on the Environment file.
type Host struct {
	OS       string    `toml:"operating_system"`
	Version  string    `toml:"go_version"`
	Hostname string    `toml:"host"`
	User     user.User `toml:"user"`
}

// Dependencies struct defines the go vendor-spec format.
type Dependencies struct {
	// Comment is free text for human use. Example "Revision abc123 introduced
	// changes that are not backwards compatible, so leave this as def876."
	Comment string `json:"comment,omitempty" toml:"comment,omitempty"`

	// Package represents a collection of vendor packages that have been copied
	// locally. Each entry represents a single Go package.
	Package []Package `json:"package" toml:"package"`
}

// Package struct defines the format for a single package.
type Package struct {
	// Import path. Example "rsc.io/pdf".
	// go get <Path> should fetch the remote package.
	Path string `json:"path" toml:"path"`

	// Origin is an import path where it was copied from. This import path
	// may contain "vendor" segments.
	//
	// If empty or missing origin is assumed to be the same as the Path field.
	Origin string `json:"origin" toml:"origin"`

	// The revision of the package. This field must be persisted by all
	// tools, but not all tools will interpret this field.
	// The value of Revision should be a single value that can be used
	// to fetch the same or similar revision.
	// Examples: "abc104...438ade0", "v1.3.5"
	Revision string `json:"revision" toml:"revision"`

	// RevisionTime is the time the revision was created. The time should be
	// parsed and written in the "time.RFC3339" format.
	RevisionTime string `json:"revisionTime" toml:"revisionTime"`

	// Comment is free text for human use.
	Comment string `json:"comment,omitempty" toml:"comment,omitempty"`
}
