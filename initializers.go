package flow

import (
	"context"

	"github.com/munichbughunter/flow/pipeline"
	"github.com/munichbughunter/flow/pipeline/clients"
	"github.com/munichbughunter/flow/pipeline/clients/cli"
	"github.com/munichbughunter/flow/pipeline/clients/dagger"
	"github.com/munichbughunter/flow/pipeline/clients/drone"
	"github.com/munichbughunter/flow/pipeline/clients/graphviz"
)

var (
	ClientCLI      string = "cli"
	ClientDrone           = "drone"
	ClientDagger          = "dagger"
	ClientGraphviz        = "graphviz"
)

func NewDefaultCollection(opts clients.CommonOpts) *pipeline.Collection {
	p := pipeline.NewCollection()
	if err := p.AddPipelines(pipeline.New(opts.Name, DefaultPipelineID)); err != nil {
		panic(err)
	}
	p.Root = []int64{DefaultPipelineID}

	return p
}

func NewMultiCollection() *pipeline.Collection {
	return pipeline.NewCollection()
}

type InitializerFunc func(context.Context, clients.CommonOpts) (pipeline.Client, error)

// The ClientInitializers define how different RunModes initialize the Flow client
var ClientInitializers = map[string]InitializerFunc{
	ClientCLI:      cli.New,
	ClientDrone:    drone.New,
	ClientDagger:   dagger.New,
	ClientGraphviz: graphviz.New,
}

func RegisterClient(name string, initializer InitializerFunc) {
	ClientInitializers[name] = initializer
}
