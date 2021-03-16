package microsoft

import (
	"reflect"
	"testing"
)

// add two numbers with most significant members first
// without modifying original list
func addTwoNumbersSig(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}
	var a []int
	var b []int
	for p := l1; p != nil; p = p.Next {
		a = append(a, p.Val)
	}
	for p := l2; p != nil; p = p.Next {
		b = append(b, p.Val)
	}
	head := &ListNode{}
	var carry int
	for i, j := len(a)-1, len(b)-1; i >= 0 || j >= 0; {
		var c, d int
		if i >= 0 {
			c = a[i]
		}
		if j >= 0 {
			d = b[j]
		}
		v := c + d + carry
		carry = v / 10
		v %= 10
		head.Next = &ListNode{
			Val:  v,
			Next: head.Next,
		}
		i--
		j--
	}
	if carry != 0 {
		head.Next = &ListNode{
			Val:  carry,
			Next: head.Next,
		}
	}
	return head.Next
}

func Test_addTwoNumbersSig(t *testing.T) {
	tests := []struct {
		name string
		l1   *ListNode
		l2   *ListNode
		want *ListNode
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addTwoNumbersSig(tt.l1, tt.l2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addTwoNumbersSig() = %v, want %v", got, tt.want)
			}
		})
	}
}
