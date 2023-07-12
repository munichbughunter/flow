// Package flow provides the primary library / client functions, types, and methods for creating Flow pipelines.
package flow

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/munichbughunter/flow/args"
	"github.com/munichbughunter/flow/cmdutil"
	"github.com/munichbughunter/flow/errors"
	"github.com/munichbughunter/flow/pipeline"
	"github.com/munichbughunter/flow/pipeline/clients"
	"github.com/munichbughunter/flow/plog"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var ErrorCancelled = errors.New("cancelled")

const DefaultPipelineID int64 = 1

// Flow is the client that is used in every pipeline to declare the steps that make up a pipeline.
// The Flow type is not thread safe. Running any of the functions from this type concurrently may have unexpected results.
type Flow struct {
	Client     pipeline.Client
	Collection *pipeline.Collection

	// Opts are the options that are provided to the pipeline from outside sources. This includes mostly command-line arguments and environment variables
	Opts    clients.CommonOpts
	Log     logrus.FieldLogger
	Version string

	// n tracks the ID of a step so that the "flow -step=" argument will function independently of the client implementation
	// It ensures that the 11th step in a Drone generated pipeline is also the 11th step in a CLI pipeline
	n        *counter
	pipeline int64

	prevPipelines []pipeline.Pipeline
}

// Pipeline returns the current Pipeline ID used in the collection.
func (s *Flow) Pipeline() int64 {
	return s.pipeline
}

func nameOrDefault(name string) string {
	if name != "" {
		return name
	}

	return "default"
}

// When allows users to define when this pipeline is executed, especially in the remote environment.
// Users can execute the pipeline as if it was triggered from the event by supplying the `-e` or `--event` argument.
// This function will overwrite any other events that were added to the pipeline.
func (s *Flow) When(events ...pipeline.Event) {
	if err := s.Collection.AddEvents(s.pipeline, events...); err != nil {
		s.Log.WithError(err).Fatalln("Failed to add events to graph")
	}
}

// Background allows users to define steps that run in the background. In some environments this is referred to as a "Service" or "Background service".
// In many scenarios, users would like to simply use a docker image with the default command. In order to accomplish that, simply provide a step without an action.
func (s *Flow) Background(steps ...pipeline.Step) {
	if err := s.validateSteps(steps...); err != nil {
		s.Log.Fatalln(err)
	}

	for i := range steps {
		steps[i].Type = pipeline.StepTypeBackground
	}

	steps = s.setup(steps...)

	if err := s.Collection.AddSteps(s.pipeline, steps...); err != nil {
		s.Log.Fatalln(err)
	}
}

// Add allows users to define steps.
// The order in which steps are ran is defined by what they provide / require.
// Some steps do not produce anything, like for example running a suite of tests for a pass/fail result.
func (s *Flow) Add(steps ...pipeline.Step) {
	steps = s.setup(steps...)

	if err := s.runSteps(steps...); err != nil {
		s.Log.Fatalln(err)
	}
}

func (s *Flow) runSteps(steps ...pipeline.Step) error {
	if err := s.validateSteps(steps...); err != nil {
		return err
	}

	if err := s.Collection.AddSteps(s.pipeline, steps...); err != nil {
		return fmt.Errorf("error adding steps '[%s]' to collection. error: %w", strings.Join(pipeline.StepNames(steps), ", "), err)
	}

	return nil
}

func (s *Flow) Cache(action pipeline.Action, c pipeline.Cacher) pipeline.Action {
	return action
}

func (s *Flow) setup(steps ...pipeline.Step) []pipeline.Step {
	for i, step := range steps {
		// Set a default image for steps that don't provide one.
		// Most pre-made steps like `yarn`, `node`, `go` steps should provide a separate default image with those utilities installed.
		if steps[i].Image == "" {
			image := "golang:1.19"
			steps[i] = step.WithImage(image)
		}

		// Set a serial / unique identifier for this step so that we can reference it using the '-step' argument consistently.
		steps[i].ID = s.n.Next()
	}

	return steps
}

func formatError(step pipeline.Step, err error) error {
	name := step.Name
	if name == "" {
		name = fmt.Sprintf("unnamed-step-%d", step.ID)
	}

	return fmt.Errorf("[name: %s, id: %d] %w", name, step.ID, err)
}

func (s *Flow) validateSteps(steps ...pipeline.Step) error {
	for _, v := range steps {
		err := s.Client.Validate(v)
		if err == nil {
			continue
		}

		if errors.Is(err, errors.ErrorSkipValidation) {
			s.Log.Warnln(formatError(v, err).Error())
			continue
		}

		return formatError(v, err)
	}

	return nil
}

func (s *Flow) watchSignals() error {
	sig := cmdutil.WatchSignals()

	return fmt.Errorf("received OS signal: %s", sig.String())
}

// Execute is the equivalent of Done, but returns an error.
// Done should be preferred in Flow pipelines as it includes sub-process handling and logging.
func (s *Flow) Execute(ctx context.Context, collection *pipeline.Collection) error {
	// Only worry about building an entire graph if we're not running a specific step.
	if step := s.Opts.Args.Step; step == nil || (*step) == 0 {
		rootArgs := pipeline.ClientProvidedArguments
		if err := collection.BuildEdges(s.Log, rootArgs...); err != nil {
			return err
		}
	}

	if err := s.Client.Done(ctx, collection); err != nil {
		return err
	}
	return nil
}

func (s *Flow) Done() {
	var (
		ctx = context.Background()
		log = s.Log
	)

	if err := execute(ctx, s.Collection, nameOrDefault(s.Opts.Name), s.Opts, s.n, s.Execute); err != nil {
		log.WithError(err).Fatalln("error in execution")
	}
}

func parseOpts() (clients.CommonOpts, error) {
	pargs, err := args.ParseArguments(os.Args[1:])
	if err != nil {
		return clients.CommonOpts{}, fmt.Errorf("error parsing arguments. Error: %w", err)
	}

	if pargs == nil {
		return clients.CommonOpts{}, fmt.Errorf("arguments list must not be nil")
	}

	// Create standard packages based on the arguments provided.
	// This would be a good place to initialize loggers, tracers, etc
	var tracer opentracing.Tracer = &opentracing.NoopTracer{}

	logger := plog.New(pargs.LogLevel)
	jaegerCfg, err := config.FromEnv()
	if err == nil {
		// Here we ignore the closer because the jaegerTracer is the closer and we will just close that.
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

func newFlow(ctx context.Context, name string) *Flow {
	opts, err := parseOpts()
	if err != nil {
		panic(fmt.Sprintf("failed to parse arguments: %s", err.Error()))
	}

	opts.Name = name
	sw := NewClient(ctx, opts, NewDefaultCollection(opts))

	// Ensure that no matter the behavior of the initializer, we still set the version on the flow object.
	sw.Version = opts.Args.Version
	sw.pipeline = DefaultPipelineID

	return sw
}

// New creates a new Flow client which is used to create pipeline a single pipeline with many steps.
// This function will panic if the arguments in os.Args do not match what's expected.
// This function, and the type it returns, are only ran inside of a Flow pipeline, and so it is okay to treat this like it is the entrypoint of a command.
// Watching for signals, parsing command line arguments, and panics are all things that are OK in this function.
// New is used when creating a single pipeline. In order to create multiple pipelines, use the NewMulti function.
func New(name string) *Flow {
	ctx := context.Background()
	rand.Seed(time.Now().Unix())
	return newFlow(ctx, name)
}

// NewWithClient creates a new Flow object with a specific client implementation.
// This function is intended to be used in very specific environments, like in tests.
func NewWithClient(opts clients.CommonOpts, client pipeline.Client) *Flow {
	rand.Seed(time.Now().Unix())
	if opts.Args == nil {
		opts.Args = &args.PipelineArgs{}
	}

	return &Flow{
		Client:     client,
		Opts:       opts,
		Log:        opts.Log,
		Collection: NewDefaultCollection(opts),
		pipeline:   DefaultPipelineID,

		n: &counter{1},
	}
}

// NewClient creates a new Flow client based on the commonopts.
// It does not check for a non-nil "Args" field.
func NewClient(ctx context.Context, c clients.CommonOpts, collection *pipeline.Collection) *Flow {
	c.Log.Infof("Initializing Flow client '%s'", c.Args.Client)
	sw := &Flow{
		n: &counter{1},
	}

	initializer, ok := ClientInitializers[c.Args.Client]
	if !ok {
		c.Log.Fatalf("Could not initialize flow. Could not find initializer for client '%s'", c.Args.Client)
		return nil
	}
	client, err := initializer(ctx, c)
	if err != nil {
		panic(err)
	}
	sw.Client = client
	sw.Collection = collection

	sw.Opts = c
	sw.Log = c.Log

	return sw
}
