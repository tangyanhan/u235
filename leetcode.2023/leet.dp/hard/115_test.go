package hard

import "testing"

func numDistinct(s string, t string) int {
	if len(t) > len(s) {
		return 0
	}

	for itt, tc := range t {
		for its, sc := range s {
			if tc == sc {

			}
		}
	}

	return 0
}

func Test_numDistinct(t *testing.T) {
	type args struct {
		s string
		t string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"rabbit", args{"rabbbit", "rabbit"}, 3},
		{"bag", args{"babgbag", "bag"}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numDistinct(tt.args.s, tt.args.t); got != tt.want {
				t.Errorf("numDistinct() = %v, want %v", got, tt.want)
			}
		})
	}
}
