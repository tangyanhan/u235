package str

import "testing"

func TestStringSetContains(t *testing.T) {
	type args struct {
		a string
		b string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"contains", args{"abcdefg", "bdacg"}, true},
		{"contains", args{"abcdefg", "bz"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSetContains(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("StringSetContains() = %v, want %v", got, tt.want)
			}
		})
	}
}
