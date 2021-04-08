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
		{
			s:      "关关雎鸠，在河之洲。窈窕淑女，君子好逑。参差荇菜，左右流之。窈窕淑女，寤寐求之。求之不得，寤寐思服。悠哉悠哉，辗转反侧。参差荇菜，左右采之。窈窕淑女，琴瑟友之。参差荇菜，左右芼之。窈窕淑女，钟鼓乐之。",
			substr: "参差荇菜，左右流之。窈窕淑女，寤寐求之。求之不得，寤寐思服。",
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
