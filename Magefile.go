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

func Build() error {
	version, err := Version()
	if err != nil {
		return err
	}

	fmt.Println("building version", version)
	return sh.Run("go",
		"build",
		"-ldflags", fmt.Sprintf("-X main.Version=%s", version),
		"-o", "./bin/flow",
		"./cmd",
	)
}

var Default = Build
