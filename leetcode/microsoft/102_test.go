package microsoft

func levelWalk(root *TreeNode, level int, result *[][]int) {
	if root == nil {
		return
	}
	if level >= len(*result) {
		*result = append(*result, []int{root.Val})
	} else {
		(*result)[level] = append((*result)[level], root.Val)
	}
	levelWalk(root.Left, level+1, result)
	levelWalk(root.Right, level+1, result)
}

func levelOrder(root *TreeNode) [][]int {
	result := make([][]int, 0)
	levelWalk(root, 0, &result)
	return result
}
