package golang

import (
	"github.com/munichbughunter/flow"
	"github.com/munichbughunter/flow/exec"
	"github.com/munichbughunter/flow/pipeline"
)

func Test(sw *flow.Flow, pkg string) pipeline.Step {
	return pipeline.NewStep(exec.RunAction("go", "test", pkg)).
		WithImage("golang:1.20").
		Requires(pipeline.ArgumentSourceFS)
}
