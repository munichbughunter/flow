package testutil

import (
	"context"

	"github.com/munichbughunter/flow"
	"github.com/munichbughunter/flow/pipeline/clients"
	"github.com/sirupsen/logrus"
)

func NewFlow(initializer flow.InitializerFunc) *flow.Flow {
	log := logrus.New()

	opts := clients.CommonOpts{
		Log: log,
	}
	client, _ := initializer(context.Background(), opts)

	return &flow.Flow{
		Opts:       opts,
		Client:     client,
		Collection: flow.NewDefaultCollection(opts),
	}
}

func NewFlowMulti(initializer flow.InitializerFunc) *flow.Flow {
	return nil
}
