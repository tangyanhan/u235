package sensetime

import (
	"reflect"
	"testing"
)

// https://leetcode-cn.com/problems/invert-binary-tree/
func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return nil
	}
	if root.Left != nil {
		invertTree(root.Left)
	}
	if root.Right != nil {
		invertTree(root.Right)
	}
	root.Left, root.Right = root.Right, root.Left
	return root
}

func Test_invertTree(t *testing.T) {
	type args struct {
		root *TreeNode
	}
	tests := []struct {
		name string
		args args
		want *TreeNode
	}{
		{
			args: args{root: NewTreeNodeFromSlice([]int{4, 2, 7, 1, 3, 6, 9})},
			want: NewTreeNodeFromSlice([]int{4, 7, 2, 9, 6, 3, 1}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := invertTree(tt.args.root); !reflect.DeepEqual(got.midOrderTraverse(), tt.want.midOrderTraverse()) {
				t.Errorf("invertTree() = %v, want %v", got, tt.want)
			}
		})
	}
}
