package main

import (
	"context"

	"github.com/munichbughunter/flow"
	"github.com/munichbughunter/flow/golang"
	"github.com/munichbughunter/flow/pipeline"
)

// "main" defines our program pipeline.
// Every pipeline step should be instantiated using the flow client (sw).
// This allows the various client modes to work properly in different scenarios, like in a CI environment or locally.
// Logic and processing done outside of the `sw.*` family of functions may not be included in the resulting pipeline.
func main() {
	sw := flow.NewMulti()
	defer sw.Done()

	sw.Add(
		sw.New("test and build", func(sw *flow.Flow) {
			sw.Add(golang.Test(sw, "./...").WithName("test"))
		}),
	)

	sw.Add(
		sw.New("create github release", func(sw *flow.Flow) {
			sw.When(
				pipeline.GitTagEvent(pipeline.GitTagFilters{}),
			)

			sw.Add(pipeline.NamedStep("am I on a tag event?", func(ctx context.Context, opts pipeline.ActionOpts) error {
				opts.Logger.Infoln("1. I'm on a tag event.")
				opts.Logger.Infoln("2. I'm on a tag event.")
				opts.Logger.Infoln("3. I'm on a tag event.")
				return nil
			}).WithImage("alpine:latest"))
		}),
	)
}
