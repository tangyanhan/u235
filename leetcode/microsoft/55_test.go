package microsoft

import "testing"

func canJump(nums []int) bool {
	maxJump := 0
	edge := len(nums) - 1
	for i := 0; i < len(nums); i++ {
		if i <= maxJump {
			v := nums[i] + i
			if maxJump < v {
				maxJump = v
			}
			if maxJump >= edge {
				return true
			}
		}
	}
	return false
}

func Test_canJump(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		nums []int
		want bool
	}{
		{
			nums: []int{2, 3, 1, 1, 4},
			want: true,
		},
		{
			nums: []int{3, 2, 1, 0, 4},
			want: false,
		},
		{
			nums: []int{0},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := canJump(tt.nums); got != tt.want {
				t.Errorf("canJump() = %v, want %v", got, tt.want)
			}
		})
	}
}
