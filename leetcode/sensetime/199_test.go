package sensetime

import (
	"container/list"
	"reflect"
	"testing"
)

func visitLevel(l int, p *TreeNode, fn func(int, *TreeNode)) {
	if p == nil {
		return
	}
	fn(l, p)
	if p.Left != nil {
		visitLevel(l+1, p.Left, fn)
	}
	if p.Right != nil {
		visitLevel(l+1, p.Right, fn)
	}
}

func rightSideView(root *TreeNode) []int {
	var levelArr []*list.List
	visitLevel(0, root, func(l int, p *TreeNode) {
		if l >= len(levelArr) {
			levelArr = append(levelArr, list.New())
		}
		levelArr[l].PushBack(p.Val)
	})
	result := make([]int, len(levelArr))
	for i, l := range levelArr {
		result[i] = l.Back().Value.(int)
	}
	return result
}

func Test_rightSideView(t *testing.T) {
	type args struct {
		root *TreeNode
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: args{root: NewTreeNodeFromSlice([]int{1, 2, 3, NullNodeVal, 5, NullNodeVal, 4})},
			want: []int{1, 3, 4},
		},
		{
			args: args{root: NewTreeNodeFromSlice([]int{})},
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rightSideView(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("rightSideView() = %v, want %v", got, tt.want)
			}
		})
	}
}
