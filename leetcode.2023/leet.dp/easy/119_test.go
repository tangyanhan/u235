package easy_test

import (
	"reflect"
	"testing"
)

func getRow(rowIndex int) []int {
	data := make([]int, rowIndex*2+2)
	prev := data[0 : rowIndex+1]
	cur := data[rowIndex+1:]

	prev[0] = 1
	for row := 0; row <= rowIndex; row++ {
		cur[row] = 1
		for col := 1; col < row; col++ {
			cur[col] = prev[col-1] + prev[col]
		}
		prev, cur = cur, prev
	}

	return prev
}

func Test_getRow(t *testing.T) {
	tests := []struct {
		name    string
		numRows int
		want    []int
	}{
		{"single", 0, []int{1}},
		{"4", 4, []int{1, 4, 6, 4, 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRow(tt.numRows); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
