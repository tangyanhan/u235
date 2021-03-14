package microsoft

import (
	"testing"
)

func reversePart(in []byte) {
	i, j := 0, len(in)-1
	for i < j {
		in[i], in[j] = in[j], in[i]
		i++
		j--
	}
}
func reverseWords(s string) string {
	buf := []byte(s)
	reversePart(buf)
	start := -1
	result := make([]byte, 0, len(buf))
	addWord := func(end int) {
		if start == -1 {
			return
		}
		if len(result) != 0 {
			result = append(result, ' ')
		}
		reversePart(buf[start:end])
		result = append(result, buf[start:end]...)
		start = -1
	}
	for i := 0; i < len(buf); i++ {
		if buf[i] != ' ' {
			if start == -1 {
				start = i
			}
			continue
		}
		addWord(i)
	}
	if start != -1 {
		addWord(len(buf))
	}
	return string(result)
}

func Test_reverseWords(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		s    string
		want string
	}{
		{
			s:    "  hello world  ",
			want: "world hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reverseWords(tt.s); got != tt.want {
				t.Errorf("reverseWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
