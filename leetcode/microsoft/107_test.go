package microsoft

import "container/list"

func levelOrderBottom(root *TreeNode) [][]int {
	if root == nil {
		return nil
	}
	queue := list.New()
	queue.PushBack(root)
	var depth int
	result := make([][]int, 0)
	for queue.Len() != 0 {
		length := queue.Len()
		for i := 0; i < length; i++ {
			p := queue.Front()
			queue.Remove(p)
			node := p.Value.(*TreeNode)
			if len(result) <= depth {
				result = append(result, []int{node.Val})
			} else {
				result[depth] = append(result[depth], node.Val)
			}

			if node.Left != nil {
				queue.PushBack(node.Left)
			}
			if node.Right != nil {
				queue.PushBack(node.Right)
			}
		}
		depth++
	}
	low, high := 0, len(result)-1
	for low < high {
		result[low], result[high] = result[high], result[low]
		low++
		high--
	}
	return result
}
