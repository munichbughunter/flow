package commands

import "github.com/munichbughunter/flow/args"

// MustParseRunArgs parses the "run" arguments from the args slice. These options are provided by the flow command and are typically not user-specified
func MustParseArgs(pargs []string) *args.PipelineArgs {
	v, err := args.ParseArguments(pargs)
	if err != nil {
		panic(err)
	}

	return v
}
