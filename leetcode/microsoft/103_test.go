package microsoft

import (
	"reflect"
	"testing"

	"leetcode/tree"
)

func visitBinTree(root *tree.TreeNode, level int, result *[][]int) {
	if root == nil {
		return
	}
	if len(*result) < level+1 {
		*result = append(*result, []int{root.Val})
	} else {
		(*result)[level] = append((*result)[level], root.Val)
	}
	visitBinTree(root.Left, level+1, result)
	visitBinTree(root.Right, level+1, result)
}

func zigzagLevelOrder(root *tree.TreeNode) [][]int {
	result := make([][]int, 0)
	visitBinTree(root, 0, &result)
	for i := 1; i < len(result); i += 2 {
		low, high := 0, len(result[i])-1
		for low < high {
			result[i][low], result[i][high] = result[i][high], result[i][low]
			low++
			high--
		}
	}
	return result
}

func Test_zigzagLevelOrder(t *testing.T) {
	tests := []struct {
		name string
		root *tree.TreeNode
		want [][]int
	}{
		{
			root: tree.NewTreeNodeFromSlice([]int{3, 9, 20, tree.NullNodeVal, tree.NullNodeVal, 15, 7}),
			want: [][]int{{3}, {20, 9}, {15, 7}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := zigzagLevelOrder(tt.root); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("zigzagLevelOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
