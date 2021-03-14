package microsoft

import (
	"reflect"
	"testing"
)

func walk(result *[]int, matrix [][]int, i, j int) (int, int) {
	const visited = -9999
	if i >= len(matrix) || j >= len(matrix[i]) || matrix[i][j] == visited {
		return -1, -1
	}

	visit := func() {
		*result = append(*result, matrix[i][j])
		matrix[i][j] = visited
	}
	// right
	for ; j < len(matrix[i]) && matrix[i][j] != visited; j++ {
		visit()
	}
	j--
	i++
	// down
	for ; i < len(matrix) && matrix[i][j] != visited; i++ {
		visit()
	}
	i--
	j--
	// left
	for ; j >= 0 && matrix[i][j] != visited; j-- {
		visit()
	}
	j++
	i--
	// up
	for ; i >= 0 && matrix[i][j] != visited; i-- {
		visit()
	}
	i++
	j++
	return i, j
}

func spiralOrder(matrix [][]int) []int {
	result := make([]int, 0, len(matrix)*len(matrix[0]))
	i, j := walk(&result, matrix, 0, 0)
	for i != -1 {
		i, j = walk(&result, matrix, i, j)
	}
	return result
}

func Test_spiralOrder(t *testing.T) {
	type args struct {
		matrix [][]int
	}
	tests := []struct {
		name   string
		matrix [][]int
		want   []int
	}{
		{
			matrix: [][]int{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
			want:   []int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7},
		},
		{
			matrix: [][]int{{1}},
			want:   []int{1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := spiralOrder(tt.matrix); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("spiralOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
