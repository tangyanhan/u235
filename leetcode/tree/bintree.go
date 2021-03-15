package tree

import "container/list"

// TreeNode binary tree
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

const NullNodeVal = -999999

// BinTreeArray is a full bin tree in array
type BinTreeArray []int

// LeftChild left child
func (b BinTreeArray) LeftChild(i int) (idx int, val int) {
	idx = i*2 + 1
	if idx >= len(b) {
		return -1, NullNodeVal
	}
	val = b[idx]
	return
}

// RightChild rightchild
func (b BinTreeArray) RightChild(i int) (idx int, val int) {
	idx = i*2 + 2
	if idx >= len(b) {
		return -1, NullNodeVal
	}
	val = b[idx]
	return
}

func (b BinTreeArray) buildTree(i int, p *TreeNode) {
	leftIdx, val := b.LeftChild(i)
	if val != NullNodeVal {
		p.Left = &TreeNode{
			Val: val,
		}
		b.buildTree(leftIdx, p.Left)
	}
	rightIdx, val := b.RightChild(i)
	if val != NullNodeVal {
		p.Right = &TreeNode{
			Val: val,
		}
		b.buildTree(rightIdx, p.Right)
	}
}

// NewTreeNodeFromSlice create new tree from slice
func NewTreeNodeFromSlice(s []int) *TreeNode {
	if len(s) == 0 {
		return nil
	}
	root := &TreeNode{
		Val: s[0],
	}
	arrTree := BinTreeArray(s)
	arrTree.buildTree(0, root)

	return root
}

// 二叉树的遍历有三种：先序，中序，后序，这里的先中后，都是指根应该在什么时候遍历访问

// 先序遍历，“根左右”
func (t *TreeNode) PreOrderTraverse() []int {
	result := make([]int, 0)
	result = append(result, t.Val)
	if t.Left != nil {
		result = append(result, t.Left.PreOrderTraverse()...)
	}
	if t.Right != nil {
		result = append(result, t.Right.PreOrderTraverse()...)
	}
	return result
}

// 中序遍历，“左根右”
func (t *TreeNode) MidOrderTraverse() []int {
	stack := list.New()
	result := make([]int, 0)
	pushNode := func(p *TreeNode) {
		for ; p != nil; p = p.Left {
			stack.PushBack(p)
		}
	}
	pushNode(t)
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

// 后序遍历，“左右根”
func (t *TreeNode) postOrderTraverse() []int {
	result := make([]int, 0)
	if t.Left != nil {
		result = append(result, t.Left.postOrderTraverse()...)
	}
	if t.Right != nil {
		result = append(result, t.Right.postOrderTraverse()...)
	}
	result = append(result, t.Val)
	return result
}
