package pingcap

import (
	"testing"
)

func insertBits(N int, M int, i int, j int) int {
	mask := 1 << i
	for k := i; k < j; k++ {
		mask |= mask << 1
	}
	M = M << i

	return mask&M | (^mask & N)
}

func toBinStr(v int) string {
	mask := 1 << 31
	buf := make([]byte, 31)
	for i := 0; i < len(buf); i++ {
		bit := mask & v
		mask = mask >> 1
		if bit != 0 {
			buf[i] = '1'
		} else {
			buf[i] = '0'
		}
	}
	return string(buf)
}

func Test_insertBits(t *testing.T) {
	type args struct {
		N int
		M int
		i int
		j int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				N: 1024,
				M: 19,
				i: 2,
				j: 6,
			},
			want: 1100,
		},
		{
			args: args{
				N: 1143207437,
				M: 1017033,
				i: 11,
				j: 31,
			},
			want: 2082885133,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := insertBits(tt.args.N, tt.args.M, tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("insertBits() = \n%v\n%v", toBinStr(got), toBinStr(tt.want))
			}
		})
	}
}
