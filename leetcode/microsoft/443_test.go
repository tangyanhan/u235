package microsoft

import "testing"

func compress(in []byte) int {
	var written int
	var start int
	writeCompressed := func(end int) {
		if end > len(in) {
			return
		}
		num := end - start
		in[written] = in[end-1]
		written++
		if num != 1 {
			numStart := written
			for num != 0 {
				v := num % 10
				in[written] = byte(v) + '0'
				written++
				num /= 10
			}
			numEnd := written - 1
			for numStart < numEnd {
				in[numStart], in[numEnd] = in[numEnd], in[numStart]
				numStart++
				numEnd--
			}
		}
		start = end
	}
	var i int
	for i = 1; i < len(in); i++ {
		if in[i] != in[i-1] {
			writeCompressed(i)
		}
	}
	if start < i {
		writeCompressed(len(in))
	}

	return written
}

func Test_compress(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want int
	}{
		{
			in:   "abbbbbb",
			want: 3,
		},
		{
			in:   "aba",
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			in := []byte(tt.in)
			if got := compress(in); got != tt.want {
				t.Errorf("compress() = %v, want %v compressed:%s", got, tt.want, in[:got])
			}
		})
	}
}
