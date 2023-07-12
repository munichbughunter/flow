package drone

import (
	"context"

	"github.com/munichbughunter/flow/pipeline"
	"github.com/munichbughunter/flow/pipeline/clients"
)

func New(ctx context.Context, opts clients.CommonOpts) (pipeline.Client, error) {
	return &Client{
		Opts: opts,
		Log:  opts.Log,
	}, nil
}
