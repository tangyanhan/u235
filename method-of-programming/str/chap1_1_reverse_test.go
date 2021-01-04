package str

import "testing"

func TestLeftRotateString(t *testing.T) {
	type args struct {
		s string
		m int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"single", args{"a", 1}, "a"},
		{"double", args{"ab", 1}, "ba"},
		{"triple", args{"abc", 1}, "bca"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LeftRotateString(tt.args.s, tt.args.m); got != tt.want {
				t.Errorf("LeftRotateString() = %v, want %v", got, tt.want)
			}
		})
	}
}
