package flow_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/munichbughunter/flow"
	"github.com/munichbughunter/flow/args"
	"github.com/munichbughunter/flow/pipeline"
	"github.com/munichbughunter/flow/pipeline/clients"
	"github.com/munichbughunter/flow/pipeline/clients/cli"
	"github.com/munichbughunter/flow/pipeline/clients/dagger"
	"github.com/munichbughunter/flow/pipeline/clients/drone"
	"github.com/munichbughunter/flow/pipeline/dag"
	"github.com/munichbughunter/flow/plog"
	"github.com/munichbughunter/flow/state"
	"github.com/sirupsen/logrus"
)

func logger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	return logger
}

var testOpts = clients.CommonOpts{
	Name:    "test piepline",
	Version: "v0.0.0-test",
	Log:     logger(),
}

func TestNew(t *testing.T) {
	t.Run("New should return a Dagger client when provided no client flag", func(t *testing.T) {
		cliArgs := []string{}
		pargs, err := args.ParseArguments(cliArgs)
		if err != nil {
			t.Fatal(err)
		}

		opts := clients.CommonOpts{
			Log:  plog.New(logrus.DebugLevel),
			Args: pargs,
		}
		sw := flow.NewClient(context.Background(), opts, flow.NewDefaultCollection(testOpts))

		if reflect.TypeOf(sw.Client) != reflect.TypeOf(&dagger.Client{}) {
			t.Fatalf("flow.Client is '%v', not a DaggerClient", reflect.TypeOf(sw.Client))
		}

		// Because reflect feels iffy to me, also make sure that it does not equal the same type as a different client
		if reflect.TypeOf(sw.Client) == reflect.TypeOf(&drone.Client{}) {
			t.Fatalf("flow.Client is '%v', not a DaggerClient", reflect.TypeOf(&drone.Client{}))
		}
	})
	t.Run("New should return a CLIClient when provided the --client=cli flag", func(t *testing.T) {
		cliArgs := []string{"--client", "cli", "--step", "1"}
		pargs, err := args.ParseArguments(cliArgs)
		if err != nil {
			t.Fatal(err)
		}

		opts := clients.CommonOpts{
			Log:  plog.New(logrus.DebugLevel),
			Args: pargs,
		}
		sw := flow.NewClient(context.Background(), opts, flow.NewDefaultCollection(testOpts))

		if reflect.TypeOf(sw.Client) != reflect.TypeOf(&cli.Client{}) {
			t.Fatalf("flow.Client is '%v', not a CLIClient", reflect.TypeOf(sw.Client))
		}

		// Because reflect feels iffy to me, also make sure that it does not equal the same type as a different client
		if reflect.TypeOf(sw.Client) == reflect.TypeOf(&drone.Client{}) {
			t.Fatalf("flow.Client is '%v', not a CLIClient", reflect.TypeOf(&drone.Client{}))
		}
	})

	t.Run("New should return a DroneClient when provided the --client=drone flag", func(t *testing.T) {
		cliArgs := []string{"--client", "drone"}
		pargs, err := args.ParseArguments(cliArgs)
		if err != nil {
			t.Fatal(err)
		}

		opts := clients.CommonOpts{
			Log:  plog.New(logrus.DebugLevel),
			Args: pargs,
		}

		sw := flow.NewClient(context.Background(), opts, flow.NewDefaultCollection(testOpts))

		if reflect.TypeOf(sw.Client) != reflect.TypeOf(&drone.Client{}) {
			t.Fatalf("flow.Client is '%v', not a DroneClient", reflect.TypeOf(sw.Client))
		}

		// Because reflect feels iffy to me, also make sure that it does not equal the same type as a different client
		if reflect.TypeOf(sw.Client) == reflect.TypeOf(&cli.Client{}) {
			t.Fatalf("flow.Client is '%v', not a DroneClient", reflect.TypeOf(&cli.Client{}))
		}
	})
}

func TestFlowRun(t *testing.T) {
	t.Run("Using a single Run function", func(t *testing.T) {
		var (
			argA = state.NewStringArgument("a")
			argB = state.NewStringArgument("b")
		)
		// In this test case we're not providing ensurer data because we are not running 'Done'.
		client := flow.NewWithClient(testOpts, newEnsurer())
		client.Add(
			pipeline.NoOpStep.WithName("step 1"),
			pipeline.NoOpStep.WithName("step 2").Provides(argA),
			pipeline.NoOpStep.WithName("step 3").Provides(argB),
			pipeline.NoOpStep.WithName("step 4").Requires(argA, argB),
		)
		// populate the graph edges
		if err := client.Collection.BuildEdges(logrus.StandardLogger()); err != nil {
			t.Fatal(err)
		}
		n, err := client.Collection.Graph.Node(flow.DefaultPipelineID)
		if err != nil {
			t.Fatal(err)
		}

		dag.EnsureGraphEdges(t, map[int64][]int64{
			0: {1, 2, 3},
			2: {4},
			3: {4},
		}, n.Value.Graph.Edges)
	})
}

func TestBasicPipeline(t *testing.T) {
	ensurer := newEnsurer("step 1", "step 2", "step 3", "step 4", "step 5")
	var (
		argA = state.NewStringArgument("a")
		argB = state.NewStringArgument("b")
	)
	client := flow.NewWithClient(testOpts, ensurer)
	client.Add(pipeline.NoOpStep.WithName("step 1").Provides(argA))
	client.Add(pipeline.NoOpStep.WithName("step 5").Requires(argB))
	client.Add(pipeline.NoOpStep.WithName("step 2").Requires(argA).Provides(argB), pipeline.NoOpStep.WithName("step 3").Requires(argA), pipeline.NoOpStep.WithName("step 4").Requires(argA))

	if err := client.Execute(context.Background(), client.Collection); err != nil {
		t.Fatal(err)
	}
}

func TestWithEvent(t *testing.T) {
	t.Run("Once adding an event, it should be present in the collection", func(t *testing.T) {
		ens := newEnsurer()
		sw := flow.NewWithClient(testOpts, ens)

		sw.When(
			pipeline.GitTagEvent(pipeline.GitTagFilters{}),
		)

		sw.Add(pipeline.NoOpStep.WithName("step 1"))

		sw.Collection.WalkPipelines(context.Background(), func(ctx context.Context, p pipeline.Pipeline) error {
			if len(p.Events) != 1 {
				t.Fatal("Expected 1 pipeline event, but found", len(p.Events))
			}

			return nil
		})
	})
}
