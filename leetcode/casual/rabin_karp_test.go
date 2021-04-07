package casual

import (
	"strings"
	"testing"
)

func TestIndexString(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		substr string
	}{
		{
			s:      "我是中国人",
			substr: "中国",
		},
		{
			s:      "关关雎鸠，在河之洲",
			substr: "左传",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := strings.Index(tt.s, tt.substr)
			if got := IndexString(tt.s, tt.substr); got != want {
				t.Errorf("IndexString() = %v, want %v", got, want)
			}
		})
	}
}
