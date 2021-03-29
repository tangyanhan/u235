package microsoft

import (
	"reflect"
	"sort"
	"testing"
)

// 寻找欧拉通路： 找到一个出度为奇数的点u
// 选择一条边(u, v)
// 删除边(u, v)
// 对v递归该过程
// 将u压入栈
func findItinerary(tickets [][]string) []string {
	g := make(map[string][]string)
	for _, ticket := range tickets {
		g[ticket[0]] = append(g[ticket[0]], ticket[1])
	}
	for k := range g {
		sort.Strings(g[k])
	}
	result := []string{}
	var dfsVisit func(string)
	dfsVisit = func(from string) {
		for {
			dests, ok := g[from]
			if !ok || len(dests) == 0 {
				break
			}
			to := dests[0]
			// 删除边
			if len(dests) == 1 {
				delete(g, from)
			} else {
				g[from] = dests[1:]
			}
			dfsVisit(to)
		}
		result = append(result, from)
	}

	dfsVisit("JFK")

	for i := 0; i < len(result)/2; i++ {
		right := len(result) - 1 - i
		result[i], result[right] = result[right], result[i]
	}
	return result
}

func Test_findItinerary(t *testing.T) {
	tests := []struct {
		name string
		in   [][]string
		want []string
	}{
		{
			in:   [][]string{{"MUC", "LHR"}, {"JFK", "MUC"}, {"SFO", "SJC"}, {"LHR", "SFO"}},
			want: []string{"JFK", "MUC", "LHR", "SFO", "SJC"},
		},
		{
			in:   [][]string{{"JFK", "SFO"}, {"JFK", "ATL"}, {"SFO", "ATL"}, {"ATL", "JFK"}, {"ATL", "SFO"}},
			want: []string{"JFK", "ATL", "JFK", "SFO", "ATL", "SFO"},
		},
		{
			in:   [][]string{{"JFK", "KUL"}, {"JFK", "NRT"}, {"NRT", "JFK"}},
			want: []string{"JFK", "NRT", "JFK", "KUL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findItinerary(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("findItinerary() = %v, want %v", got, tt.want)
			}
		})
	}
}
