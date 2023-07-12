package fs

import (
	"context"

	"github.com/munichbughunter/flow/pipeline"
)

func Replace(file string, content string) pipeline.Action {
	return func(context.Context, pipeline.ActionOpts) error {
		return nil
	}
}

func ReplaceString(file string, content string) pipeline.Action {
	return func(context.Context, pipeline.ActionOpts) error {
		return nil
	}
}
