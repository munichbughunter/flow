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

func (g *Graph[T]) AddEdge(from, to int64) error {
	var fromNode, toNode *Node[T]

	for i, v := range g.Nodes {
		if v.ID == from {
			fromNode = &g.Nodes[i]
		}
		if v.ID == to {
			toNode = &g.Nodes[i]
		}
		if fromNode != nil && toNode != nil {
			break
		}
	}

	if fromNode == nil {
		return fmt.Errorf("%w. id: %d", ErrorNotFound, from)
	}
	if toNode == nil {
		return fmt.Errorf("%w. id: %d", ErrorNotFound, to)
	}

	edges := g.Edges[from]
	g.Edges[from] = append(edges, Edge[T]{
		From: fromNode,
		To:   toNode,
	})

	return nil
}

func (g *Graph[T]) Node(id int64) (*Node[T], error) {
	for i, v := range g.Nodes {
		if v.ID == id {
			return &g.Nodes[i], nil
		}
	}

	return nil, fmt.Errorf("id: %d. error: %w", id, ErrorNotFound)
}

func (g *Graph[T]) NodeList(id ...int64) ([]*Node[T], error) {
	nodes := make([]*Node[T], len(id))
	for i, v := range id {
		node, err := g.Node(v)
		if err != nil {
			return nil, nil, fmt.Errorf("id: %d. error: %w", v, err)
		}

		nodes[i] = node
	}

	return nodes, nil
}
