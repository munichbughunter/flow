package main

import "flow"

var Pipelines = []flow.Pipeline{
	PipelineDependencies,
	PipelineBuild,
	PipelineTest,
	PipelinePublish,
}

// "main" defines our program pipeline.
// Every pipeline step should be instantiated using the flow client (sw).
// This allows the various clients to work properly in different scenarios, like in a CI environment or locally.
// Logic and processing done outside of the `sw.*` family of functions may not be included in the resulting pipeline.
func main() {
	sw := flow.NewMulti()
	defer sw.Done()

	sw.AddPipelines(Pipelines...)
}
