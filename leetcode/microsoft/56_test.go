package microsoft

import (
	"reflect"
	"sort"
	"testing"
)

func merge(intervals [][]int) [][]int {
	sort.Slice(intervals, func(i, j int) bool {
		if intervals[i][0] != intervals[j][0] {
			return intervals[i][0] < intervals[j][0]
		}
		return intervals[i][1] < intervals[j][1]
	})
	results := make([][]int, 0)
	var last []int
	for _, vec := range intervals {
		if last != nil {
			if vec[0] >= last[0] && vec[0] <= last[1] {
				if vec[1] >= last[1] {
					last[1] = vec[1]
				}
				continue
			}
		}
		results = append(results, vec)
		last = vec
	}
	return results
}

func Test_merge(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want [][]int
	}{
		{
			in:   [][]int{{1, 4}, {2, 3}},
			want: [][]int{{1, 4}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := merge(tt.in); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("merge() = %v, want %v", got, tt.want)
			}
		})
	}
}
