package easy

// Definition for a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

const NullLeaf = -1

func NewNodeFromSlice(data []int, index int) *TreeNode {
	if index >= len(data) || data[index] == NullLeaf {
		return nil
	}

	root := new(TreeNode)
	root.Val = data[index]

	root.Left = NewNodeFromSlice(data, 2*index+1)
	root.Right = NewNodeFromSlice(data, 2*index+2)

	return root
}
