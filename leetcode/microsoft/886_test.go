package microsoft

import (
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
	set := NewUnionSet(N*2 + 1)
	for _, pair := range dislikes {
		px := set.Find(pair[0])
		py := set.Find(pair[1])
		if px == py {
			return false
		}
		disX := set.Find(pair[0] + N)
		disY := set.Find(pair[1] + N)

		set.Join(px, disY)
		set.Join(py, disX)
	}
	return true
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
