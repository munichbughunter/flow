package main

import (
	"context"

	"github.com/munichbughunter/flow"
	"github.com/munichbughunter/flow/exec"
	"github.com/munichbughunter/flow/pipeline"
)

func echo(ctx context.Context, opts pipeline.ActionOpts) error {
	return exec.RunCommandWithOpts(ctx, exec.RunOpts{
		Name:   "/bin/sh",
		Args:   []string{"-c", `sleep 10; echo "hello ?"`},
		Stdout: opts.Stdout,
		Stderr: opts.Stderr,
	})
}

func StepEcho() pipeline.Step {
	return pipeline.NewStep(echo).WithImage("ubuntu:latest")
}

func main() {
	sw := scribe.New("test-pipeline")
	defer sw.Done()

	sw.Add(StepEcho())
}
