package microsoft

import "testing"

func countComponents(n int, edges [][]int) int {
	set := NewUnionSet(n)
	for _, edge := range edges {
		set.Join(edge[0], edge[1])
	}
	return set.Count
}

func Test_countComponents(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges [][]int
		want  int
	}{
		{
			n:     5,
			edges: [][]int{{0, 1}, {1, 2}, {3, 4}},
			want:  2,
		},
		{
			n:     5,
			edges: [][]int{{0, 1}, {1, 2}, {2, 3}, {3, 4}},
			want:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := countComponents(tt.n, tt.edges); got != tt.want {
				t.Errorf("countComponents() = %v, want %v", got, tt.want)
			}
		})
	}
}
