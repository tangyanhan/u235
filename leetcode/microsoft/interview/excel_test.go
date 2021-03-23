package interview

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

type Node struct {
	Cell string
	Exps []*Node
}

type Graph map[string]*Node

func (g Graph) HasSCC() bool {
	visited := make(map[string]int)
	var visit func(node *Node) bool
	var currCount int
	visit = func(node *Node) bool {
		if node == nil || len(node.Exps) == 0 {
			return false
		}
		cnt, ok := visited[node.Cell]
		if ok {
			return cnt == currCount
		}

		visited[node.Cell] = currCount
		for _, subNode := range node.Exps {
			if visit(subNode) {
				return true
			}
		}
		return false
	}

	for _, node := range g {
		if visit(node) {
			return true
		}
		currCount++
	}
	return false
}

func evaluate(in []string) bool {
	regCell := regexp.MustCompile("([A-Z]{2}[0-9]{2})")
	graph := make(Graph)
	getOrAddNode := func(cell string) *Node {
		node, ok := graph[cell]
		if !ok {
			node = &Node{
				Cell: cell,
			}
			graph[cell] = node
		}
		return node
	}
	// build graph
	for _, s := range in {
		parts := strings.Split(s, " = ")
		cell := parts[0]
		node := getOrAddNode(cell)
		exps := regCell.FindAllStringSubmatch(parts[1], -1)
		if len(exps) == 0 {
			continue
		}
		fmt.Println("in:", parts[1], "Submatches:", exps)
		for _, exp := range exps {
			// found another graph node
			expNode := getOrAddNode(exp[0])
			// 自环
			if node.Cell == expNode.Cell {
				return false
			}
			node.Exps = append(node.Exps, expNode)
		}
	}
	// graph now: node-> []subnodes
	return !graph.HasSCC()
}

func Test_evaluate(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want bool
	}{
		{
			in: []string{
				"AA00 = 10",
				"AA01 = AA00 + AB00",
				"AB00 = 15",
			},
			want: true,
		},
		{
			in: []string{
				"AA00 = 10",
				"AB00 = (AA00 + AA01) * 15 + AA00",
				"AA01 = 20 + AB00",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := evaluate(tt.in); got != tt.want {
				t.Errorf("evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
