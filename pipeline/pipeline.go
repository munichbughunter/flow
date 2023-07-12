package pipeline

import "errors"

var (
	ErrorNoStepProvider    = errors.New("no step in the graph provides a required argument")
	ErrorAmbiguousProvider = errors.New("more than one step provides the same argument(s)")
)

type Pipeline struct {
	ID        int64
	Name      string
	Graph     *dag.Graph[Step]
	Providers map[state.Argument]int64
	Root      []int64
	Events    []Event
	Type      PipelineType

	RequiredArgs state.Arguments
	ProvidedArgs state.Arguments
}
