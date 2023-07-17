package main

import (
	"fmt"
	"os"
)

func InitFlow(lang Language) {
	err := os.MkdirAll("./flow", 0755)
	if err != nil {
		panic(err)
	}

	var created bool
	files := []string{"flow", "build", "publish", "deploy", "test"}
	for _, file := range files {
		filename := "flow/" + file + lang.GetFileExtension()
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			// Creating folder and files!
			f, err := os.Create(filename)
			if err != nil {
				panic(err)
			}
			f.Close()
			created = true
		} else {
			created = false
		}
	}

	GenerateGitlabYml(lang)

	if created {
		fmt.Println("Flow project initialized")
	} else {
		fmt.Println("Flow project already exists! Nothing to create!")
	}
}
