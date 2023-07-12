package golang

import (
	"context"

	"github.com/munichbughunter/flow/golang/x"
	"github.com/munichbughunter/flow/pipeline"
)

func BuildStep(pkg, output string, args, env []string) pipeline.Step {
	return pipeline.NewStep(func(ctx context.Context, opts pipeline.ActionOpts) error {
		return x.RunBuild(ctx, x.BuildOpts{
			Pkg:    pkg,
			Output: output,
			Stdout: opts.Stdout,
			Stderr: opts.Stderr,
			Env:    env,
			Args:   args,
		})
	})
}

func BuildAction(pkg, output string, args, env []string) pipeline.Action {
	return func(ctx context.Context, opts pipeline.ActionOpts) error {
		opts.Logger.Infoln("args: ", args)
		return x.RunBuild(ctx, x.BuildOpts{
			Pkg:    pkg,
			Output: output,
			Stdout: opts.Stdout,
			Stderr: opts.Stderr,
			Env:    env,
			Args:   args,
		})
	}
}
