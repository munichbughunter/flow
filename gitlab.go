// gitlab.go
package main

import (
	"fmt"
	"os"
)

func GenerateGitlabYml(lang Language) {
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

	if _, err := os.Stat(".gitlab-ci.yml"); os.IsNotExist(err) {
		err := os.WriteFile(".gitlab-ci.yml", []byte(gitlabYml), 0666)
		if err != nil {
			panic(err)
		}
	}

	err := os.WriteFile(".gitlab-ci.yml", []byte(gitlabYml), 0666)
	if err != nil {
		panic(err)
	}
	fmt.Println("gitlab-ci.yml file created successfully")
}
