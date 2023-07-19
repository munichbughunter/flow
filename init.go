package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func InitFlow(lang Language) {
	templateDir := "templates"
	flowDir := "flow"

	err := os.MkdirAll("./flow", 0755)
	if err != nil {
		panic(err)
	}

	prefix := strings.TrimPrefix(lang.GetFileExtension(), ".")
	srcDir := filepath.Join(templateDir, prefix)
	dstDir := flowDir

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			filename := filepath.Base(path)
			newFilename := strings.Replace(filename, ".temp", "."+prefix, 1)
			newPath := filepath.Join(dstDir, newFilename)

			input, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			err = ioutil.WriteFile(newPath, input, 0755)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

	GenerateGitlabYml(lang)
}
