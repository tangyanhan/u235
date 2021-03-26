package microsoft

import (
	"fmt"
	"testing"
)

type intQueue struct {
	data  []int
	start int
	end   int
	size  int
}

func newIntQueue(capacity int) *intQueue {
	return &intQueue{
		data: make([]int, capacity),
	}
}

func (q *intQueue) push(v int) {
	q.data[q.end] = v
	q.end++
	q.size++
	if q.end == len(q.data) {
		q.end = 0
	}
}

func (q *intQueue) pop() int {
	v := q.data[q.start]
	q.start++
	q.size--
	if q.start == len(q.data) {
		q.start = 0
	}
	return v
}

func (q *intQueue) length() int {
	return q.size
}

func possibleBipartition(N int, dislikes [][]int) bool {
	var g [2001][2001]bool
	var inDegree [2001]int
	vertexMap := make(map[int]bool)
	for _, pair := range dislikes {
		fmt.Println(pair[0], "<->", pair[1])
		g[pair[0]][pair[1]] = true
		g[pair[1]][pair[0]] = true
		inDegree[pair[0]]++
		inDegree[pair[1]]++
		vertexMap[pair[0]] = true
		vertexMap[pair[1]] = true
	}
	fmt.Println(inDegree)
	queue := newIntQueue(2000)
	for i, v := range inDegree {
		if v == 1 {
			queue.push(i)
		}
	}
	var numVisited int
	for queue.length() != 0 {
		from := queue.pop()
		fmt.Println("Visited:", from)
		numVisited++
		for to := range vertexMap {
			if g[from][to] {
				inDegree[to]--
				inDegree[from]--
				if inDegree[to] == 1 {
					queue.push(to)
				}
			}
		}
	}
	return len(vertexMap) == numVisited
}

func Test_possibleBipartition(t *testing.T) {
	tests := []struct {
		name     string
		N        int
		dislikes [][]int
		want     bool
	}{
		{
			N:        3,
			dislikes: [][]int{{1, 2}, {2, 3}},
			want:     true,
		},
		{
			N:        3,
			dislikes: [][]int{{1, 2}, {2, 3}, {1, 3}},
			want:     false,
		},
		{
			N:        10,
			dislikes: [][]int{{4, 7}, {4, 8}, {5, 6}, {1, 6}, {3, 7}, {2, 5}, {5, 8}, {1, 2}, {4, 9}, {6, 10}, {8, 10}, {3, 6}, {2, 10}, {9, 10}, {3, 9}, {2, 3}, {1, 9}, {4, 6}, {5, 7}, {3, 8}, {1, 8}, {1, 7}, {2, 4}},
			want:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := possibleBipartition(tt.N, tt.dislikes); got != tt.want {
				t.Errorf("possibleBipartition() = %v, want %v", got, tt.want)
			}
		})
	}
}
