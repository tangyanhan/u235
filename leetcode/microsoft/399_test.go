package microsoft

import (
	"math"
	"testing"
)

func calcEquation(equations [][]string, values []float64, queries [][]string) []float64 {
	type exp struct {
		dest  string
		value float64
	}
	g := make(map[string][]exp)
	for i, pair := range equations {
		from, to := pair[0], pair[1]
		g[from] = append(g[from], exp{dest: to, value: values[i]})
		// values[i] > 0.0，所以不用处理除0异常
		g[to] = append(g[to], exp{dest: from, value: 1 / values[i]})
	}
	result := make([]float64, len(queries))
	var target string
	var idx int
	var dfsVisit func(visited map[string]bool, curr string, value float64)
	dfsVisit = func(visited map[string]bool, curr string, value float64) {
		if curr == target {
			result[idx] = value
			return
		}
		visited[curr] = true
		for _, edge := range g[curr] {
			if visited[edge.dest] {
				continue
			}
			if edge.dest == target {
				result[idx] = value * edge.value
				return
			}
			dfsVisit(visited, edge.dest, value*edge.value)
		}
	}
	var query []string
	for idx, query = range queries {
		result[idx] = -1.0
		from, to := query[0], query[1]
		// 点不存在
		if _, ok := g[from]; !ok {
			continue
		}
		if _, ok := g[to]; !ok {
			continue
		}
		target = to
		visited := make(map[string]bool)
		dfsVisit(visited, from, 1.0)
	}
	return result
}

func Test_calcEquation(t *testing.T) {
	type args struct {
		equations [][]string
		values    []float64
		queries   [][]string
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{
		{
			args: args{
				equations: [][]string{{"a", "b"}, {"b", "c"}},
				values:    []float64{2.0, 3.0},
				queries:   [][]string{{"a", "c"}, {"b", "a"}, {"a", "e"}, {"a", "a"}, {"x", "x"}},
			},
			want: []float64{6.0, 0.5, -1.0, 1.0, -1.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcEquation(tt.args.equations, tt.args.values, tt.args.queries)
			for i, v := range tt.want {
				if math.Abs(v-got[i]) > 0.001 {
					t.Errorf("calcEquation() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
