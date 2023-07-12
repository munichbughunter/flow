package main

import (
	"context"

	"github.com/munichbughunter/flow"
	"github.com/munichbughunter/flow/fs"
	gitx "github.com/munichbughunter/flow/git/x"
	"github.com/munichbughunter/flow/golang"
	"github.com/munichbughunter/flow/makefile"
	"github.com/munichbughunter/flow/pipeline"
	"github.com/munichbughunter/flow/state"
	"github.com/munichbughunter/flow/yarn"
)

var (
	ArgumentTestResult = state.NewBoolArgument("test-results")
)

func writeVersion(sw *flow.Flow) pipeline.Step {
	action := func(ctx context.Context, opts pipeline.ActionOpts) error {

		// equivalent of `git describe --tags --dirty --always`
		version, err := gitx.Describe(ctx, ".", true, true, true)
		if err != nil {
			return err
		}

		// write the version string in the `.version` file.
		return fs.ReplaceString(".version", version)(ctx, opts)
	}

	return pipeline.NewStep(action)
}

func installDependencies(sw *flow.Flow) {
	sw.Add(
		pipeline.NamedStep("install frontend dependencies", sw.Cache(
			yarn.InstallAction(),
			fs.Cache("node_modules", fs.FileHasChanged("yarn.lock")),
		)),
		pipeline.NamedStep("install backend dependencies", sw.Cache(
			golang.ModDownload(),
			fs.Cache("$GOPATH/pkg", fs.FileHasChanged("go.sum")),
		)),
	)
}

func testPipeline(sw *flow.Flow) {
	installDependencies(sw)

	sw.Add(
		golang.Test(sw, "./...").WithName("test backend"),
		pipeline.NamedStep("test frontend", makefile.Target("test-frontend")),
	)
}

func publishPipeline(sw *flow.Flow) {
	sw.When(
		pipeline.GitCommitEvent(pipeline.GitCommitFilters{
			Branch: pipeline.StringFilter("main"),
		}),
		pipeline.GitTagEvent(pipeline.GitTagFilters{
			Name: pipeline.GlobFilter("v*"),
		}),
	)

	installDependencies(sw)

	sw.Add(
		pipeline.NamedStep("compile backend", makefile.Target("build")),
		pipeline.NamedStep("compile frontend", makefile.Target("package")),
	)

	sw.Add(
		pipeline.NamedStep("publish", makefile.Target("publish")).Requires(state.NewSecretArgument("gcp-publish-key")),
	)
}

func codeqlPipeline(sw *flow.Flow) {
	sw.Add(
		pipeline.NoOpStep.WithName("codeql"),
		pipeline.NoOpStep.WithName("notify-slack"),
	)
}

// "main" defines our program pipeline.
// Every pipeline step should be instantiated using the flow client (sw).
// This allows the various clients to work properly in different scenarios, like in a CI environment or locally.
// Logic and processing done outside of the `sw.*` family of functions may not be included in the resulting pipeline.
func main() {
	sw := flow.NewMulti()
	defer sw.Done()

	sw.Add(
		sw.New("code quality check", codeqlPipeline),
		sw.New("test", testPipeline).Provides(ArgumentTestResult),
		sw.New("publish", publishPipeline).Requires(ArgumentTestResult),
	)
}
