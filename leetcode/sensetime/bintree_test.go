package sensetime

import (
	"reflect"
	"testing"
)

func TestTraverse(t *testing.T) {
	tree := NewTreeNodeFromSlice([]int{0, 1, 2, 3, 4, 5, 6, NullNodeVal, 8, 9, NullNodeVal})
	t.Log(tree.preOrderTraverse())
	t.Log(tree.midOrderTraverse())
	t.Log(tree.postOrderTraverse())
}

func TestTreeNode_preOrderTraverse(t *testing.T) {
	tests := []struct {
		name string
		t    *TreeNode
		want []int
	}{
		{
			t:    NewTreeNodeFromSlice([]int{0, 1, 2, 3, 4, 5, 6, NullNodeVal, 8, 9, NullNodeVal}),
			want: []int{0, 1, 3, 8, 4, 9, 2, 5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.preOrderTraverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TreeNode.preOrderTraverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTreeNode_midOrderTraverse(t *testing.T) {
	tests := []struct {
		name string
		t    *TreeNode
		want []int
	}{
		{
			t:    NewTreeNodeFromSlice([]int{0, 1, 2, 3, 4, 5, 6, NullNodeVal, 8, 9, NullNodeVal}),
			want: []int{3, 8, 1, 9, 4, 0, 5, 2, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.midOrderTraverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TreeNode.midOrderTraverse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTreeNode_postOrderTraverse(t *testing.T) {
	tests := []struct {
		name string
		t    *TreeNode
		want []int
	}{
		{
			t:    NewTreeNodeFromSlice([]int{0, 1, 2, 3, 4, 5, 6, NullNodeVal, 8, 9, NullNodeVal}),
			want: []int{8, 3, 9, 4, 1, 5, 6, 2, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.postOrderTraverse(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TreeNode.postOrderTraverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
