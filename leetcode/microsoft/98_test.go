package microsoft

import (
	"container/list"
	"fmt"
	"math"
	"testing"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func helper(low, high int64, root *TreeNode) bool {
	if root == nil {
		return true
	}
	if int64(root.Val) <= low || int64(root.Val) >= high {
		return false
	}
	return helper(low, int64(root.Val), root.Left) && helper(int64(root.Val), high, root.Right)
}

func isValidBSTRecursive(root *TreeNode) bool {
	if root == nil {
		return true
	}

	return helper(math.MinInt64, math.MaxInt64, root)
}

func isValidBST(root *TreeNode) bool {
	stack := list.New()
	minValue := int64(math.MinInt64)
	for stack.Len() != 0 || root != nil {
		for root != nil {
			fmt.Println("Pushed:", root.Val)
			stack.PushBack(root)
			root = root.Left
		}
		root = stack.Back().Value.(*TreeNode)
		fmt.Println("Popped=", root.Val, "min=", minValue)
		if int64(root.Val) <= minValue {
			return false
		}
		stack.Remove(stack.Back())
		minValue = int64(root.Val)
		root = root.Right
	}
	return true
}

func Test_isValidBST(t *testing.T) {
	tests := []struct {
		name string
		root *TreeNode
		want bool
	}{
		{
			root: &TreeNode{
				Val: 2,
				Left: &TreeNode{
					Val: 1,
				},
				Right: &TreeNode{
					Val: 3,
				},
			},
			want: true,
		},
		{
			root: &TreeNode{
				Val: 5,
				Left: &TreeNode{
					Val: 1,
				},
				Right: &TreeNode{
					Val: 4,
					Left: &TreeNode{
						Val: 3,
					},
					Right: &TreeNode{
						Val: 6,
					},
				},
			},
			want: false,
		},
		{
			root: &TreeNode{
				Val: math.MaxInt32,
			},
			want: true,
		},
		{
			root: &TreeNode{
				Val: math.MaxInt32,
				Left: &TreeNode{
					Val: math.MaxInt32,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidBST(tt.root); got != tt.want {
				t.Errorf("isValidBST() = %v, want %v", got, tt.want)
			}
		})
	}
}
