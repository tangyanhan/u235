package microsoft

import (
	"container/list"
	"fmt"
	"testing"
)

func longestIncreasingPath(matrix [][]int) int {
	// TODO: provide correct solution
	// idx -> i * n + j
	n := len(matrix[0])
	g := make(map[int][]int)
	inDegree := make(map[int]int)
	outDegree := make(map[int]int)
	hasEdge := func(curr int, i, j int) (int, bool) {
		if i < 0 || j < 0 || i >= len(matrix) || j >= len(matrix[0]) {
			return -1, false
		}
		if matrix[i][j] > curr {
			fmt.Println("Added edge:", curr, "->", matrix[i][j], "i=", i, "j=", j)
			idx := i*n + j
			inDegree[idx]++
			return idx, true
		}
		return -1, false
	}

	var base int
	for i := 0; i < len(matrix); i++ {
		base = i * n
		for j := 0; j < len(matrix[i]); j++ {
			var s []int
			curr := matrix[i][j]
			fmt.Println("For num:", curr, "pos:", i, j)
			idx := base + j
			// up
			if v, ok := hasEdge(curr, i-1, j); ok {
				s = append(s, v)
			}
			// down
			if v, ok := hasEdge(curr, i+1, j); ok {
				s = append(s, v)
			}
			// left
			if v, ok := hasEdge(curr, i, j-1); ok {
				s = append(s, v)
			}
			// right
			if v, ok := hasEdge(curr, i, j+1); ok {
				s = append(s, v)
			}
			g[idx] = s
			outDegree[idx] = len(s)
		}
	}
	queue := list.New()
	for k, degree := range inDegree {
		fmt.Println("degree:", degree, matrix[k/n][k%n])
		if degree == 0 {
			queue.PushBack(k)
		}
	}
	depth := 1
	for queue.Len() != 0 {
		length := queue.Len()
		for i := 0; i < length; i++ {
			pop := queue.Front()
			queue.Remove(pop)
			from := pop.Value.(int)
			fmt.Println("Visit:", from)
			for _, to := range g[from] {
				inDegree[to]--
				if inDegree[to] == 0 {
					queue.PushBack(to)
				}
			}
		}
		depth++
	}
	return depth
}

func Test_longestIncreasingPath(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		want   int
	}{
		{
			matrix: [][]int{{3, 4, 5}, {3, 2, 6}, {2, 2, 1}},
			want:   4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := longestIncreasingPath(tt.matrix); got != tt.want {
				t.Errorf("longestIncreasingPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
