package main

type Language interface {
	GetDockerImage() string
	GetBuildScript() string
	GetFileExtension() string
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

func (p Python) GetFileExtension() string {
	return ".py"
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

func (ts TypeScript) GetFileExtension() string {
	return ".mts"
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

func (golang Golang) GetFileExtension() string {
	return ".go"
}
