package main

import (
	"context"

	"github.com/munichbughunter/flow"
	"github.com/munichbughunter/flow/pipeline"
	"github.com/munichbughunter/flow/pipeline/clients"
	"github.com/munichbughunter/flow/state"
	"github.com/sirupsen/logrus"
)

type MyClient struct {
	Log logrus.FieldLogger
}

func (c *MyClient) Validate(step pipeline.Step) error {
	return nil
}

func (c *MyClient) Provides() []state.Argument {
	return nil
}

func (c *MyClient) Done(ctx context.Context, w *pipeline.Collection) error {
	return w.WalkPipelines(ctx, func(ctx context.Context, p pipeline.Pipeline) error {
		c.Log.Infoln("pipeline:", p.Name)
		return w.WalkSteps(ctx, p.ID, func(ctx context.Context, step pipeline.Step) error {
			c.Log.Infoln("step:", step.Name)
			return nil
		})
	})
}

func init() {
	flow.RegisterClient("my-custom-client", func(ctx context.Context, opts clients.CommonOpts) (pipeline.Client, error) {
		return &MyClient{
			Log: opts.Log,
		}, nil
	})
}

func main() {
	sw := flow.New("custom-client")
	defer sw.Done()

	sw.Add(
		pipeline.NoOpStep.WithName("step 1"),
		pipeline.NoOpStep.WithName("step 2"),
	)
	sw.Add(
		pipeline.NoOpStep.WithName("step 3"),
		pipeline.NoOpStep.WithName("step 4"),
	)
}
