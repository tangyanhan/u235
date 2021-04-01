package casual

import (
	"testing"
)

func clumsy(N int) int {
	var ans int
	var last int
	for n, op := N, 0; n >= 1; n-- {
		switch op {
		case 0:
			last = n
		case 1:
			last *= n
		case 2:
			last /= n
		case 3:
			last += n
		case 4:
			ans += last
			last = 0 - n
		}
		op++
		if op == 5 {
			op = 1
		}
	}
	ans += last
	return ans
}

func Test_clumsy(t *testing.T) {
	tests := []struct {
		name string
		N    int
	}{
		{
			N: 10,
		},
		{
			N: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := 0
			var last int
			for n, op := tt.N, 0; n >= 1; n-- {
				switch op {
				case 0:
					last = n
				case 1:
					last *= n
				case 2:
					last /= n
				case 3:
					last += n
				case 4:
					want += last
					last = 0 - n
				}
				t.Log(n, "last=", last, "want=", want, "op=", op)
				op++
				if op == 5 {
					op = 1
				}
			}
			t.Log("last=", last, "want=", want)
			want += last

			if got := clumsy(tt.N); got != want {
				t.Errorf("clumsy() = %v, want %v", got, want)
			}
		})
	}
}
