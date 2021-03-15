package microsoft

import (
	"reflect"
	"testing"
)

func postOrderWalk(root, p, q *TreeNode) (*TreeNode, bool) {
	if root == nil {
		return nil, false
	}
	var found bool
	if root == p || root == q {
		found = true
	}
	ans, foundLeft := postOrderWalk(root.Left, p, q)
	if ans != nil {
		return ans, false
	}
	ans, foundRight := postOrderWalk(root.Right, p, q)
	if ans != nil {
		return ans, false
	}
	if (foundLeft && foundRight) || (found && (foundLeft || foundRight)) {
		return root, false
	}

	return nil, foundLeft || foundRight || found
}

func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	ans, _ := postOrderWalk(root, p, q)
	return ans
}

func Test_lowestCommonAncestor(t *testing.T) {
	type args struct {
		root *TreeNode
		p    *TreeNode
		q    *TreeNode
	}
	tests := []struct {
		name string
		args args
		want *TreeNode
	}{
		{
			args: func() args {
				root := &TreeNode{Val: 3}
				n5 := &TreeNode{Val: 5}
				n1 := &TreeNode{Val: 1}
				root.Left, root.Right = n5, n1
				return args{
					root: root,
					p:    n5,
					q:    n1,
				}
			}(),
		},
		{
			args: func() args {
				root := &TreeNode{
					Val: 3,
					Left: &TreeNode{
						Val: 5,
						Left: &TreeNode{
							Val: 6,
						},
						Right: &TreeNode{
							Val: 2,
							Left: &TreeNode{
								Val: 7,
							},
							Right: &TreeNode{
								Val: 4,
							},
						},
					},
					Right: &TreeNode{
						Val: 1,
						Left: &TreeNode{
							Val: 0,
						},
						Right: &TreeNode{
							Val: 8,
						},
					},
				}
				return args{
					root: root,
					p:    root.Left,
					q:    root.Left.Right.Right,
				}
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lowestCommonAncestor(tt.args.root, tt.args.p, tt.args.q); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lowestCommonAncestor() = %v, want %v", got, tt.want)
			}
		})
	}
}
