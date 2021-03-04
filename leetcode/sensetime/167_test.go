package sensetime

import (
	"reflect"
	"testing"
)

var sumResult [2]int

func twoSum(numbers []int, target int) []int {
	low, high := 0, len(numbers)-1
	for low < high {
		remain := target - numbers[low]
		if remain == numbers[high] {
			sumResult[0], sumResult[1] = low+1, high+1
			return sumResult[:]
		}
		if remain < numbers[high] {
			high--
		} else {
			low++
		}
	}
	return nil
}

func Test_twoSum(t *testing.T) {
	type args struct {
		numbers []int
		target  int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: args{
				numbers: []int{2, 7, 11, 25},
				target:  9,
			},
			want: []int{1, 2},
		},
		{
			args: args{
				numbers: []int{2, 3, 4},
				target:  6,
			},
			want: []int{1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := twoSum(tt.args.numbers, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("twoSum() in=%v = %v, want %v", tt.args.numbers, got, tt.want)
			}
		})
	}
}
