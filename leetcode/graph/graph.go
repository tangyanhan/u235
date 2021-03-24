package graph

import (
	"container/list"
	"fmt"
	"math"
)

// Graph instance
type Graph struct {
	Vertexs        map[int]*Vertex
	NegativeWeight int
}

// InitSingleSource initialize single source
func (g *Graph) InitSingleSource(source *Vertex) {
	for _, v := range g.Vertexs {
		v.UpperBound = math.MaxInt32
		v.Pi = nil
	}
	source.UpperBound = 0
}

// Vertex node
type Vertex struct {
	Value      int
	UpperBound int
	Pi         *Vertex
	Edges      *Edge
	LastEdge   *Edge
}

// Relax u->v
func (v *Vertex) Relax(edge *Edge) {
	if edge.To.UpperBound > v.UpperBound+edge.Weight {
		edge.To.UpperBound = v.UpperBound + edge.Weight
		edge.To.Pi = v
	}
}

// Edge edges
type Edge struct {
	To     *Vertex
	Weight int

	Next *Edge
	Prev *Edge
}

// VisitVertexFunc defines how to visit node in DFS/BFS. depth starts from 0.
type VisitVertexFunc func(vertex *Vertex, depth int) (shouldContinue bool)

// New alloc new graph
func New() *Graph {
	return &Graph{
		Vertexs: make(map[int]*Vertex),
	}
}

// GetOrAddVertex add node on non-exists
func (g *Graph) GetOrAddVertex(v int) *Vertex {
	node, ok := g.Vertexs[v]
	if !ok {
		node = &Vertex{
			Value: v,
		}
		g.Vertexs[v] = node
	}
	return node
}

// AddEdge add edge
func (g *Graph) AddEdge(from, to, weight int) {

	fromVertex := g.GetOrAddVertex(from)
	toVertex := g.GetOrAddVertex(to)
	for edge := fromVertex.Edges; edge != nil; edge = edge.Next {
		if edge.To == toVertex {
			if edge.Weight < 0 {
				g.NegativeWeight -= edge.Weight
			}
			if weight < 0 {
				g.NegativeWeight += weight
			}
			edge.Weight = weight
			return
		}
	}
	if weight < 0 {
		g.NegativeWeight += weight
	}
	edge := &Edge{
		To:     toVertex,
		Weight: weight,
	}
	if fromVertex.Edges == nil {
		fromVertex.Edges = edge
		fromVertex.LastEdge = edge
		return
	}
	fromVertex.LastEdge.Next = edge
	edge.Prev = fromVertex.LastEdge
}

// DFS depth first
func (g *Graph) DFS(start *Vertex, fn VisitVertexFunc) {
	visited := make(map[int]bool)
	var visitVertex func(node *Vertex, depth int) bool
	visitVertex = func(node *Vertex, depth int) bool {
		if node == nil {
			return true
		}
		if visited[node.Value] {
			// previous rounds
			return true
		}
		visited[node.Value] = true
		if !fn(node, depth) {
			return false
		}
		for edge := node.Edges; edge != nil; edge = edge.Next {
			if !visitVertex(edge.To, depth+1) {
				return false
			}
		}
		return true
	}
	visitVertex(start, 0)
}

// BFS broad first
func (g *Graph) BFS(start *Vertex, fn VisitVertexFunc) {
	queue := list.New()
	visited := make(map[int]bool)
	queue.PushBack(start)
	visited[start.Value] = true
	var depth int
	for queue.Len() != 0 {
		length := queue.Len()
		for i := 0; i < length; i++ {
			p := queue.Front()
			node := p.Value.(*Vertex)
			if !fn(node, depth) {
				fmt.Println("Step=", depth)
				return
			}
			queue.Remove(p)
			for edge := node.Edges; edge != nil; edge = edge.Next {
				if !visited[edge.To.Value] {
					queue.PushBack(edge.To)
					visited[edge.To.Value] = true
				}
			}
		}
		depth++
	}
}
