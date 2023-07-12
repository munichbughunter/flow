package flow

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Flow struct {
	Client     pipeline.Client
	Collection *pipeline.Collection

	Opts    clients.CommonOpts
	Log     logrus.FieldLogger
	Version string

	n             *counter
	pipeline      int64
	prevPipelines []pipeline.Pipeline
}

func New(name string) *Flow {
	ctx := context.Background()
	rand.Seed(time.Now().Unix())
	return newFlow(ctx, name)
}

func newFlow(ctx context.Context, name string) *Flow {
	opts, err := parseOpts()
	if err != nil {
		panic(fmt.Sprintf("failed to parse arguments: %s", err.Error()))
	}
	opts.Name = name
	fw := NewClient(ctx, opts, NewDefaultCollection(opts))

	fw.Version = opts.Args.Version
	fw.pipeline = DefaultPipelineID

	return fw
}

func parseOpts() (clients.CommonOpts, error) {
	pargs, err := args.ParseArguments(os.Args[1:])
	if err != nil {
		return clients.CommonOpts{}, fmt.Errorf("error parsing arguments. Error: %w", err)
	}

	if pargs == nil {
		return clients.CommonOpts{}, fmt.Errorf("arguments list must not be nil")
	}

	var tracer opentracing.Tracer = &opentracing.NoopTracer{}

	logger := plog.New(pargs.LogLevel)
	jaegerCfg, err := config.FromEnv()
	if err == nil {
		jaegerTracer, _, err := jaegerCfg.NewTracer(config.Logger(jaeger.StdLogger))
		if err == nil {
			logger.Debugln("Initialized jaeger tracer")
			tracer = jaegerTracer
		} else {
			logger.Debugln("Could not initialize jaeger tracer; using no-op tracer; Error:", err.Error())
		}
	}

	return clients.CommonOpts{
		Version: pargs.Version,
		Output:  os.Stdout,
		Args:    pargs,
		Log:     logger,
		Tracer:  tracer,
	}, nil
}
