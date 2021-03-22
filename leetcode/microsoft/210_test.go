package microsoft

import (
	"reflect"
	"testing"
)

type Graph [][]int

func (g Graph) AddEdge(from, to int) {
	g[from][to] = 1
}

func (g Graph) RemoveEdge(from, to int) {
	g[from][to] = 0
}

// GraphNode holds value
type GraphNode struct {
	Val     int
	Parents []*GraphNode
	Childs  []*GraphNode
}

func findOrder(numCourses int, prerequisites [][]int) []int {
	graph := make(map[int]*GraphNode)
	badRoots := make(map[int]bool)
	for i := 0; i < numCourses; i++ {
		graph[i] = &GraphNode{
			Val: i,
		}
	}

	for _, pair := range prerequisites {
		dependNode := graph[pair[1]]
		node := graph[pair[0]]
		dependNode.Childs = append(dependNode.Childs, node)
		node.Parents = append(node.Parents, dependNode)
		badRoots[pair[0]] = true
	}
	visited := make([]bool, numCourses)
	result := make([]int, 0, numCourses)
	var newResultAdded bool
	var visitNode func(*GraphNode)
	// DFS on node
	visitNode = func(node *GraphNode) {
		if visited[node.Val] {
			return
		}
		// if this one has parents not visited yet, skip for the time being
		if len(node.Parents) != 0 {
			for _, parentNode := range node.Parents {
				if !visited[parentNode.Val] {
					return
				}
			}
		}
		visited[node.Val] = true
		result = append(result, node.Val)
		newResultAdded = true

		for _, child := range node.Childs {
			visitNode(child)
		}
	}
	for len(result) != numCourses {
		newResultAdded = false
		for i := 0; i < numCourses; i++ {
			visitNode(graph[i])
		}
		if !newResultAdded {
			break
		}
	}
	if len(result) != numCourses {
		return nil
	}
	return result
}

func Test_findOrder(t *testing.T) {
	type args struct {
		numCourses    int
		prerequisites [][]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: args{
				prerequisites: [][]int{{1, 0}, {2, 0}, {3, 1}, {3, 2}},
				numCourses:    4,
			},
			want: []int{0, 2, 1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findOrder(tt.args.numCourses, tt.args.prerequisites); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
