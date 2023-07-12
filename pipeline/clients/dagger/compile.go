package dagger

import (
	"context"
	"path/filepath"

	"dagger.io/dagger"
	"github.com/munichbughunter/flow/pipelineutil"
)

func CompilePipeline(ctx context.Context, d *dagger.Client, name, src, gomod, pipeline string) (*dagger.Directory, error) {
	var (
		dir     = d.Host().Directory(src)
		builder = d.Container().From("golang:1.19").WithMountedDirectory("/src", dir)
	)

	path, err := filepath.Rel(src, gomod)
	if err != nil {
		return nil, err
	}
	cmd := pipelineutil.GoBuild(ctx, pipelineutil.GoBuildOpts{
		Pipeline: pipeline,
		Module:   path,
		Output:   "/opt/flow/pipeline",
		LDFlags:  `-extldflags "-static"`,
	})

	builder = builder.WithEnvVariable("GOOS", "linux")
	builder = builder.WithEnvVariable("GOARCH", "amd64")
	builder = builder.WithEnvVariable("CGO_ENABLED", "0")
	// Set the pipeline name to prevent cache collisions.
	// Some pipelines with the exact same name and path will sometimes reuse the compiled pipeline from the cache.
	// In those scenarios, until we find a more permanent fix, it's best to just change the name.
	builder = builder.WithEnvVariable("PIPELINE_NAME", name)
	builder = builder.WithWorkdir("/src")

	builder = builder.Exec(dagger.ContainerExecOpts{
		Args: cmd.Args,
	})

	return builder.Directory("/opt/flow"), nil
}
