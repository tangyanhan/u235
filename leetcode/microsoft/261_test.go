package microsoft

import (
	"container/list"
	"fmt"
	"testing"
)

func validTreeBFS(n int, edges [][]int) bool {
	g := make(map[int][]int)
	color := make([]int, n+1)
	for _, edge := range edges {
		from, to := edge[0], edge[1]
		g[from] = append(g[from], to)
		g[to] = append(g[to], from)
	}

	var numVisited int
	queue := list.New()
	queue.PushBack(0)

	for queue.Len() != 0 {
		length := queue.Len()
		for i := 0; i < length; i++ {
			pop := queue.Front()
			queue.Remove(pop)
			from := pop.Value.(int)
			if color[from] != 0 {
				continue
			}
			color[from] = 1
			numVisited++
			for _, to := range g[from] {
				if color[to] == 0 {
					fmt.Println("Pushed:", from, " <->", to)
					queue.PushBack(to)
				}
			}
			color[from] = 2
		}
	}
	return numVisited == n
}

func validTree(n int, edges [][]int) bool {
	g := make(map[int][]int)
	color := make([]int, n+1)
	for _, edge := range edges {
		from, to := edge[0], edge[1]
		g[from] = append(g[from], to)
		g[to] = append(g[to], from)
	}

	var numVisited int
	var dfsVisit func(int, int) bool
	dfsVisit = func(curr, prev int) bool {
		if color[curr] == 2 {
			return false
		}
		if color[curr] == 1 {
			return true
		}
		numVisited++
		color[curr] = 1
		for _, to := range g[curr] {
			if prev != to {
				if dfsVisit(to, curr) {
					return true
				}
			}
		}
		color[curr] = 2
		return false
	}

	if dfsVisit(0, -1) || numVisited != n {
		return false
	}
	return true
}

func Test_validTree(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges [][]int
		want  bool
	}{
		{
			n:     5,
			edges: [][]int{{0, 1}, {1, 2}, {2, 3}, {1, 3}, {1, 4}},
			want:  false,
		},
		{
			n:     4,
			edges: [][]int{{0, 1}, {2, 3}},
			want:  false,
		},
		{
			n:     4,
			edges: [][]int{{0, 1}, {0, 2}, {1, 2}},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validTree(tt.n, tt.edges); got != tt.want {
				t.Errorf("validTree() = %v, want %v", got, tt.want)
			}
			if got := validTreeBFS(tt.n, tt.edges); got != tt.want {
				t.Errorf("validTreeBFS() = %v, want %v", got, tt.want)
			}
		})
	}
}
