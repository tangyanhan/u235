package str

import (
	"testing"
)

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

func TestReverseSentence(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "normal",
			in:   "He is a student.",
			want: "student. a is He",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ReverseSentence(tt.in); got != tt.want {
				t.Errorf("ReverseSentence() = %v, want %v", got, tt.want)
			}
		})
	}
}
