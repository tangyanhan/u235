package microsoft

import (
	"container/list"
)

func minDepth(root *TreeNode) int {
	queue := list.New()
	depth := 1
	if root != nil {
		queue.PushBack(root)
	}
	for queue.Len() != 0 {
		length := queue.Len()
		for i := 0; i < length; i++ {
			p := queue.Front()
			node := p.Value.(*TreeNode)
			if node.Left == nil && node.Right == nil {
				return depth
			}
			queue.Remove(p)
			if node.Left != nil {
				queue.PushBack(node.Left)
			}
			if node.Right != nil {
				queue.PushBack(node.Right)
			}
		}
		depth++
	}
	return 0
}
