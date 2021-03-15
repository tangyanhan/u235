package microsoft

import (
	"leetcode/tree"
	"reflect"
	"testing"
)

func buildTree(preorder []int, inorder []int) *tree.TreeNode {
	if len(preorder) == 0 {
		return nil
	}
	root := &tree.TreeNode{
		Val: preorder[0],
	}
	var rootIdx int
	for ; inorder[rootIdx] != root.Val; rootIdx++ {
	}
	root.Left = buildTree(preorder[1:rootIdx+1], inorder[:rootIdx])
	root.Right = buildTree(preorder[rootIdx+1:], inorder[rootIdx+1:])
	return root
}

func TestBuildTree(t *testing.T) {
	root := tree.NewTreeNodeFromSlice([]int{5, 4, 6, tree.NullNodeVal, tree.NullNodeVal, 3, 7})
	preorder := root.PreOrderTraverse()
	inorder := root.MidOrderTraverse()
	t.Log("Preorder:", preorder, "Inorder:", inorder)
	newTree := buildTree(preorder, inorder)
	npre := newTree.PreOrderTraverse()
	nin := newTree.MidOrderTraverse()
	if !reflect.DeepEqual(npre, preorder) || !reflect.DeepEqual(inorder, nin) {
		t.Fatal(npre, preorder, nin, inorder)
	}
}
