package main

import (
	"context"
	"time"

	"github.com/munichbughunter/flow"
	"github.com/munichbughunter/flow/pipeline"
	"github.com/munichbughunter/flow/state"
)

var (
	ArgumentTestResultBackend  = state.NewBoolArgument("backend-test-result")
	ArgumentTestResultFrontend = state.NewBoolArgument("frontend-test-result")
)

func actionTestFrontend(ctx context.Context, opts pipeline.ActionOpts) error {
	opts.Logger.Infoln("Testing frontend...")
	time.Sleep(time.Second * 1)
	opts.Logger.Infoln("Done testing frontend")
	// make test-frontend
	// assume it passed...
	return opts.State.SetBool(ctx, ArgumentTestResultFrontend, true)
}

func actionTestBackend(ctx context.Context, opts pipeline.ActionOpts) error {
	opts.Logger.Infoln("Testing backend...")
	time.Sleep(time.Second * 1)
	opts.Logger.Infoln("Done testing backend")
	// go test ./...
	// assume it passed...
	return opts.State.SetBool(ctx, ArgumentTestResultBackend, true)
}

var stepTestBackend = pipeline.NamedStep("test backend", actionTestBackend).
	Provides(ArgumentTestResultBackend).
	Requires(ArgumentGoDependencies)

var stepTestFrontend = pipeline.NamedStep("test frontend", actionTestFrontend).
	Provides(ArgumentTestResultFrontend).
	Requires(ArgumentNodeDependencies)

var PipelineTest = flow.Pipeline{
	Name: "test",
	Steps: []pipeline.Step{
		stepTestBackend,
		stepTestFrontend,
	},
	Requires: []state.Argument{ArgumentNodeDependencies, ArgumentGoDependencies},
	Provides: []state.Argument{ArgumentTestResultFrontend, ArgumentTestResultBackend},
}
