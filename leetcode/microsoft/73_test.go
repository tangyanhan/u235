package microsoft

import (
	"reflect"
	"testing"
)

func setZeroes(matrix [][]int) {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return
	}
	var shouldCleanFirstRow bool
	var shouldCleanFirstCol bool

	for col := 0; col < len(matrix[0]); col++ {
		if matrix[0][col] == 0 {
			shouldCleanFirstRow = true
			break
		}
	}
	for row := 0; row < len(matrix); row++ {
		if matrix[row][0] == 0 {
			shouldCleanFirstCol = true
			break
		}
	}

	cleanRow := func(row int) {
		for col := 0; col < len(matrix[row]); col++ {
			matrix[row][col] = 0
		}
	}
	cleanCol := func(col int) {
		for row := 0; row < len(matrix); row++ {
			matrix[row][col] = 0
		}
	}

	// mark if rows should be cleaned at col[0]
	for row := 1; row < len(matrix); row++ {
		for col := 1; col < len(matrix[0]); col++ {
			if matrix[row][col] == 0 {
				matrix[row][0] = 0
				matrix[0][col] = 0
			}
		}
	}
	// clean marked rows
	for row := 1; row < len(matrix); row++ {
		if matrix[row][0] == 0 {
			cleanRow(row)
		}
	}
	// clean marked cols
	for col := 1; col < len(matrix[0]); col++ {
		if matrix[0][col] == 0 {
			cleanCol(col)
		}
	}

	if shouldCleanFirstRow {
		cleanRow(0)
	}
	if shouldCleanFirstCol {
		cleanCol(0)
	}
}

func Test_setZeroes(t *testing.T) {
	tests := []struct {
		name   string
		matrix [][]int
		want   [][]int
	}{
		{
			matrix: [][]int{{0, 1, 2, 0}, {3, 4, 5, 2}, {1, 3, 1, 5}},
			want:   [][]int{{0, 0, 0, 0}, {0, 4, 5, 0}, {0, 3, 1, 0}},
		},
		{
			matrix: [][]int{{0}, {1}},
			want:   [][]int{{0}, {0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setZeroes(tt.matrix)
			if !reflect.DeepEqual(tt.matrix, tt.want) {
				t.Fatalf("Want: %v\n Got:%v", tt.want, tt.matrix)
			}
		})
	}
}
