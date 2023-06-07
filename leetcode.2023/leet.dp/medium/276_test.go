package medium_test

import "testing"

func numWays(n int, k int) int {
	return 0
}

func Test_numWays(t *testing.T) {
	tests := []struct {
		name string
		n    int
		k    int
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numWays(tt.n, tt.k); got != tt.want {
				t.Errorf("numWays() = %v, want %v", got, tt.want)
			}
		})
	}
}
