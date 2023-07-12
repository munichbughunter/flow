package testutil

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/munichbughunter/flow/args"
	"github.com/munichbughunter/flow/cmd/commands"
	"github.com/munichbughunter/flow/pipeline"
)

func RunPipeline(ctx context.Context, t *testing.T, path string, stdout io.Writer, stderr io.Writer, args *args.PipelineArgs) {
	stderrBuf := bytes.NewBuffer(nil)
	stdoutBuf := bytes.NewBuffer(nil)
	t.Log("Running pipeline with args", args)
	cmd := commands.Run(ctx, &commands.RunOpts{
		Path:   path,
		Stdout: io.MultiWriter(stdout, stdoutBuf),
		Stderr: io.MultiWriter(stderr, stderrBuf),
		Args:   args,
	})

	if err := cmd.Run(); err != nil {
		t.Fatalf("Error running pipeline. Error: '%s'\nStdout: '%s'\nStderr: '%s'\n", err, stdoutBuf.String(), stderrBuf.String())
	}
}

// NewTestStep creates a new TestStep that emits data into the channel 'b' when the action is ran
func NewTestStep(b chan bool) pipeline.Step {
	return pipeline.Step{
		Name: "test",
		Action: func(context.Context, pipeline.ActionOpts) error {
			b <- true
			return nil
		},
	}
}
