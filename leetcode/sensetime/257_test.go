package sensetime

import (
	"container/list"
	"reflect"
	"strconv"
	"testing"
)

var paths []string

// https://leetcode-cn.com/problems/binary-tree-paths
func binaryTreePaths(root *TreeNode) []string {
	paths = []string{}
	stack := list.New()
	type nodeVal struct {
		Node *TreeNode
		Str  string
	}
	pushNodes := func(root *TreeNode) {
		for ; root != nil; root = root.Left {
			if stack.Back() != nil {
				stack.PushBack(&nodeVal{
					Node: root,
					Str:  stack.Back().Value.(*nodeVal).Str + "->" + strconv.Itoa(root.Val),
				})
			} else {
				stack.PushBack(&nodeVal{
					Node: root,
					Str:  strconv.Itoa(root.Val),
				})
			}
		}
	}
	pushNodes(root)
	for stack.Len() != 0 {
		p := stack.Back()
		v := p.Value.(*nodeVal)
		stack.Remove(p)
		if v.Node.Left == nil && v.Node.Right == nil {
			paths = append(paths, v.Str)
		}
		if v.Node.Right != nil {
			pushNodes(v.Node.Right)
		}
	}
	return paths
}

func constructPaths(root *TreeNode, path string) {
	if root != nil {
		pathSB := path
		pathSB += strconv.Itoa(root.Val)
		if root.Left == nil && root.Right == nil {
			paths = append(paths, pathSB)
		} else {
			pathSB += "->"
			constructPaths(root.Left, pathSB)
			constructPaths(root.Right, pathSB)
		}
	}
}

func Test_binaryTreePaths(t *testing.T) {
	type args struct {
		root *TreeNode
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			args: args{root: NewTreeNodeFromSlice([]int{1, 2, 3, NullNodeVal, 5})},
			want: []string{"1->2->5", "1->3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := binaryTreePaths(tt.args.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("binaryTreePaths() = %v, want %v", got, tt.want)
			}
		})
	}
}
