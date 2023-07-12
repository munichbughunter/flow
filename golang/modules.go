package golang

import (
	"context"

	"github.com/munichbughunter/flow/pipeline"
)

func ModDownload() pipeline.Action {
	return func(context.Context, pipeline.ActionOpts) error {
		return nil
	}
}
