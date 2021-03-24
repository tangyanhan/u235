package graph

import (
	"fmt"
	"testing"
)

func TestDFS(t *testing.T) {

}

func TestGraph_DFS(t *testing.T) {
	type testCase struct {
		name  string
		edges [][]int
		fn    VisitVertexFunc
		check func() string
	}
	tests := []testCase{
		func() testCase {
			c := testCase{
				name: "traverse",
				edges: [][]int{
					{0, 1}, {0, 2}, {0, 3},
					{2, 3}, {2, 4}, {2, 5},
					{6, 6},
				},
			}
			visited := make(map[int]bool)
			c.fn = func(node *Vertex, depth int) bool {
				visited[node.Value] = true
				fmt.Println("Visited node:", node.Value, "Depth=", depth)
				return true
			}
			c.check = func() string {
				for i := 0; i < 6; i++ {
					if !visited[i] {
						return fmt.Sprint("node not visited:", i)
					}
				}
				return ""
			}

			return c
		}(),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := New()
			for _, edge := range tt.edges {
				g.AddEdge(edge[0], edge[1], 0)
			}
			g.DFS(g.GetOrAddVertex(0), tt.fn)
			if v := tt.check(); v != "" {
				t.Fatal(v)
			}
			g.BFS(g.GetOrAddVertex(0), func(node *Vertex, depth int) bool {
				fmt.Println("BFS - ", node.Value, "depth=", depth)
				return node.Value != 5
			})
		})
	}
}
