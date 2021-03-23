package graph

import (
	"container/list"
	"fmt"
)

// Graph instance
type Graph struct {
	Nodes map[int]*Node
}

// Node node
type Node struct {
	Value  int
	Childs []*Node
}

// VisitNodeFunc defines how to visit node in DFS/BFS. depth starts from 0.
type VisitNodeFunc func(node *Node, depth int) (shouldContinue bool)

// New alloc new graph
func New() *Graph {
	return &Graph{
		Nodes: make(map[int]*Node),
	}
}

// GetOrAddNode add node on non-exists
func (g *Graph) GetOrAddNode(v int) *Node {
	node, ok := g.Nodes[v]
	if !ok {
		node = &Node{
			Value: v,
		}
		g.Nodes[v] = node
	}
	return node
}

// AddEdge add edge
func (g *Graph) AddEdge(from, to int) {
	fromNode := g.GetOrAddNode(from)
	toNode := g.GetOrAddNode(to)
	for _, child := range fromNode.Childs {
		if child.Value == to {
			return
		}
	}
	fromNode.Childs = append(fromNode.Childs, toNode)
}

// DFS depth first
func (g *Graph) DFS(start *Node, fn VisitNodeFunc) {
	visited := make(map[int]bool)
	var visitNode func(node *Node, depth int) bool
	visitNode = func(node *Node, depth int) bool {
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
		for _, child := range node.Childs {
			if !visitNode(child, depth+1) {
				return false
			}
		}
		return true
	}
	visitNode(start, 0)
}

// BFS broad first
func (g *Graph) BFS(start *Node, fn VisitNodeFunc) {
	queue := list.New()
	visited := make(map[int]bool)
	queue.PushBack(start)
	visited[start.Value] = true
	var depth int
	for queue.Len() != 0 {
		length := queue.Len()
		for i := 0; i < length; i++ {
			p := queue.Front()
			node := p.Value.(*Node)
			if !fn(node, depth) {
				fmt.Println("Step=", depth)
				return
			}
			queue.Remove(p)
			for _, child := range node.Childs {
				if !visited[child.Value] {
					queue.PushBack(child)
					visited[child.Value] = true
				}
			}
		}
		depth++
	}
}
