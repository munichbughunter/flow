package main

type Language interface {
	GetDockerImage() string
	GetBuildScript() string
}

type Python struct{}

func (p Python) GetDockerImage() string {
	return "python:3.11"
}

func (p Python) GetBuildScript() string {
	return `script:
    - python -m pip install --upgrade pip
    - pip install dagger-io
    - python ./flow/flow.py`
}

type TypeScript struct{}

func (ts TypeScript) GetDockerImage() string {
	return "node:20"
}

func (ts TypeScript) GetBuildScript() string {
	return `script:
    - npm ci
    - npm run build
    - npm test`
}

type Golang struct{}

func (golang Golang) GetDockerImage() string {
	return "golang:1.20"
}

func (golang Golang) GetBuildScript() string {
	return `script:
    - go get dagger-io
    - go run ./flow/flow.go`
}
