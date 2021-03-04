package set_bilibili

import "testing"

var fibMap = make([]int, 101)

const fibCeil = 1e9 + 7

func init() {
	fibMap[0] = 0
	fibMap[1] = 1
	for i := 2; i < len(fibMap); i++ {
		fibMap[i] = fibMap[i-1] + fibMap[i-2]
		fibMap[i] %= fibCeil
	}
}

func fib(n int) int {
	return fibMap[n]
}

func Test_fib(t *testing.T) {
	t.Log(fibMap)
	tests := []struct {
		name string
		in   int
		want int
	}{
		{
			in:   2,
			want: 1,
		},
		{
			in:   5,
			want: 5,
		},
		{
			in:   45,
			want: 134903163,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fib(tt.in); got != tt.want {
				t.Errorf("fib() = %v, want %v, map value=%v", got, tt.want, fibMap[tt.in])
			}
		})
	}
}
