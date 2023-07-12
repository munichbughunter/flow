package args

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type PipelineArgs struct {
	Client  string
	Path    string
	Version string
	Step    *int64
}
