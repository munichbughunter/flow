package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	createGitlabFile := flag.Bool("create-gitlab-file", false, "Generate gitlab.yml file")
	ts := flag.Bool("ts", false, "Use TS")
	py := flag.Bool("py", false, "Use PY")
	golang := flag.Bool("go", false, "Use GO")
	// pyFlag := flag.Bool("py", false, "Use PY template")
	// tsFlag := flag.Bool("ts", false, "Use TS template")
	flag.Parse()

	// if *createGitLabFile {
	// 	if *pyFlag {
	// 		generatePythonPipeline()
	// 	} else if *tsFlag {
	// 		generateTSPipeline()
	// 	} else {
	// 		log.Fatal("Please specify a valid flag (--python or --ts) to create a GitLab CI file.")
	// 	}
	// } else {
	// 	// Handle other commands or actions
	// }

	if *createGitlabFile {
		var lang Language
		if *ts {
			lang = TypeScript{}
		} else if *py {
			lang = Python{}
		} else if *golang {
			lang = Golang{}
		} else {
			// Error / Logging?
		}

		gitlabYml := fmt.Sprintf(
			`.docker:
  image: %s
  services:
    - docker:${DOCKER_VERSION}-dind
  variables:
    DOCKER_HOST: tcp://docker:2376
    DOCKER_TLS_VERIFY: '1'
    DOCKER_TLS_CERTDIR: '/certs'
    DOCKER_CERT_PATH: '/certs/client'
    DOCKER_DRIVER: overlay2
    DOCKER_VERSION: '20.10.16'
.dagger:
  extends: [.docker]
  before_script:
    - apk add docker-cli
build:
  extends: [.dagger]
  %s`, lang.GetDockerImage(), lang.GetBuildScript())

		err := os.WriteFile(".gitlab-ci.yml", []byte(gitlabYml), 0666)
		if err != nil {
			panic(err)
		}
		fmt.Println("gitlab.yml file created successfully")
	}
}
