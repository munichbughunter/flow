package dag

import (
	"errors"
	"fmt"
)

var (
	ErrorDuplicateID = errors.New("node with ID already exists")
	ErrorNotFound    = errors.New("node with ID not found")
	ErrorNoVisitFunc = errors.New("no visitfunc provided")

	ErrorBreak = errors.New("break will stop the depth first search without an error")
)

type Node[T any] struct {
	ID    int64
	Value T
}

type Edge[T any] struct {
	From *Node[T]
	To   *Node[T]
}

type Graph[T any] struct {
	Nodes []Node[T]
	Edges map[int64][]Edge[T]

	visited map[int64]bool
}

func (g *Graph[T]) AddNode(id int64, v T) error {
	node := Node[T]{
		ID:    id,
		Value: v,
	}

	for _, v := range g.Nodes {
		if v.ID == node.ID {
			return fmt.Errorf("%w. id: %d", ErrorDuplicateID, node.ID)
		}
	}

	g.Nodes = append(g.Nodes, node)
	return nil
}
