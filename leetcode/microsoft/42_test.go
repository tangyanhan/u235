package microsoft

import "testing"

func trap(height []int) int {
	// TODO:
	return 0
}

func Test_trap(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
	}{
		{
			in:   []int{4, 2, 0, 3, 2, 5},
			want: 9,
		},
		{
			in:   []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1},
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trap(tt.in); got != tt.want {
				t.Errorf("trap() = %v, want %v", got, tt.want)
			}
		})
	}
}
