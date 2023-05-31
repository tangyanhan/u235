package easy_test

import (
	"reflect"
	"testing"
)

func generate(numRows int) [][]int {
	result := make([][]int, numRows)
	for row := 0; row < numRows; row++ {
		result[row] = make([]int, row+1)
		result[row][0] = 1
		result[row][row] = 1
		for col := 1; col < row; col++ {
			result[row][col] = result[row-1][col] + result[row-1][col-1]
		}
	}

	return result
}

func Test_generate(t *testing.T) {
	tests := []struct {
		name    string
		numRows int
		want    [][]int
	}{
		{"single", 1, [][]int{{1}}},
		{"5", 5, [][]int{{1}, {1, 1}, {1, 2, 1}, {1, 3, 3, 1}, {1, 4, 6, 4, 1}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generate(tt.numRows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generate() = %v, want %v", got, tt.want)
			}
		})
	}
}
