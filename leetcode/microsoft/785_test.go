package microsoft

import "testing"

func isBipartite(graph [][]int) bool {
	set := NewUnionSet(2 * len(graph))
	for from, dests := range graph {
		for _, to := range dests {
			pFrom := set.Find(from)
			pTo := set.Find(to)
			if pFrom == pTo {
				return false
			}
			pNoFrom := set.Find(from + len(graph))
			pNoTo := set.Find(to + len(graph))
			set.Join(pNoFrom, to)
			set.Join(pNoTo, from)
		}
	}
	return true
}

func Test_isBipartite(t *testing.T) {
	tests := []struct {
		name  string
		graph [][]int
		want  bool
	}{
		{
			graph: [][]int{{1, 2, 3}, {0, 2}, {0, 1, 3}, {0, 2}},
			want:  false,
		},
		{
			graph: [][]int{{1, 3}, {0, 2}, {1, 3}, {0, 2}},
			want:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isBipartite(tt.graph); got != tt.want {
				t.Errorf("isBipartite() = %v, want %v", got, tt.want)
			}
		})
	}
}
