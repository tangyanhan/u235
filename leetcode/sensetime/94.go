package sensetime

import "container/list"

// 二叉树的中序遍历，压栈
// https://leetcode-cn.com/problems/binary-tree-inorder-traversal/
func inorderTraversal(root *TreeNode) []int {
	stack := list.New()
	result := make([]int, 0)
	pushNode := func(p *TreeNode) {
		for ; p != nil; p = p.Left {
			stack.PushBack(p)
		}
	}
	pushNode(root)
	for stack.Len() > 0 {
		p := stack.Back()
		node := p.Value.(*TreeNode)
		result = append(result, node.Val)
		stack.Remove(p)
		if node.Right != nil {
			pushNode(node.Right)
		}
	}

	return result
}
