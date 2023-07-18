//go:build mage
// +build mage

package main

import (
	"fmt"

	"github.com/magefile/mage/sh"
)

func Version() (string, error) {
	return sh.Output("git", "describe", "--tags", "--dirty", "--always")
}

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	version, err := Version()
	if err != nil {
		return err
	}

	fmt.Println("Building Flow Version", version)
	return sh.Run("go",
		"build",
		"-ldflags", fmt.Sprintf("-X main.Version=%s", version),
		"-o", "./bin/flow",
	)
}

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build
