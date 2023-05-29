package easy

import (
	"log"
	"testing"
)

/**
 * Definition for a binary tree node.
 * type TreeNode struct {
 *     Val int
 *     Left *TreeNode
 *     Right *TreeNode
 * }
 */
func hasPathSum(root *TreeNode, targetSum int) bool {
	if root == nil {
		return false
	}
	targetSum -= root.Val

	if targetSum == 0 && root.Left == nil && root.Right == nil {
		return true
	}

	return hasPathSum(root.Left, targetSum) || hasPathSum(root.Right, targetSum)
}

func hasPathSumComplex(root *TreeNode, targetSum int) bool {
	var walk func(node *TreeNode, sum int) bool

	walk = func(node *TreeNode, sum int) bool {
		if node == nil {
			return false
		}
		// Leaf
		if node.Left == nil && node.Right == nil {
			log.Printf("Leaf value at %d=%d \n", node.Val, node.Val+sum)
			return sum+node.Val == targetSum
		}

		if node.Left != nil {
			if walk(node.Left, sum+node.Val) {
				return true
			}
		}
		if node.Right != nil {
			if walk(node.Right, sum+node.Val) {
				return true
			}
		}
		return false
	}

	return walk(root, 0)
}

func Test_HasPathSum(t *testing.T) {
	tests := []struct {
		name   string
		data   []int
		target int
		expect bool
	}{
		{"exp-1", []int{5, 4, 8, 11, NullLeaf, 13, 4, 7, 2, NullLeaf, NullLeaf, NullLeaf, 1}, 22, true},
		{"no-value", []int{1, 2, 3}, 5, false},
		{"nil tree", []int{}, 0, false},
	}

	for _, td := range tests {
		t.Run(td.name, func(t *testing.T) {
			tree := NewNodeFromSlice(td.data, 0)
			got := hasPathSum(tree, td.target)
			if got != td.expect {
				t.Fatalf("Expect: %t Got: %t", td.expect, got)
			}
		})
	}
}
