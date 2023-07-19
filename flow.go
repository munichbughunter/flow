package main

import (
	"flag"
)

func main() {
	flowPath := flag.String("./flow", "", "Path to Flow Folder")
	init := flag.Bool("init", false, "Initialize flow")
	createGitlabFile := flag.Bool("create-gitlab-file", false, "Generate gitlab.yml file")
	help := flag.Bool("help", false, "Flow help")

	ts := flag.Bool("ts", false, "Use TS")
	py := flag.Bool("py", false, "Use PY")
	golang := flag.Bool("go", false, "Use GO")
	flag.Parse()

	var lang Language

	if *ts {
		lang = TypeScript{}
	} else if *py {
		lang = Python{}
	} else if *golang {
		lang = Golang{}
	} else {
		// panic("Language not specified! Please specify your language py, ts or go")
		// Log out that no language was specified... ?
	}

	if *flowPath != "" {
		// RunFlowPipeline(*flowPath)
	} else if *init {
		InitFlow(lang)
	} else if *createGitlabFile {
		GenerateGitlabYml(lang)
	} else if *help {
		PrintHelp()
	} else {

	}
}
