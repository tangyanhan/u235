package easy

type ListNode struct {
	Val  int
	Next *ListNode
}

func NewListFromInts(s []int) *ListNode {
	if len(s) == 0 {
		return nil
	}
	head := new(ListNode)
	var p *ListNode
	for _, v := range s {
		if p == nil {
			p = head
			p.Val = v
		} else {
			p.Next = &ListNode{
				Val: v,
			}
			p = p.Next
		}
	}
	return head
}

func ListToInts(l *ListNode) []int {
	var s []int
	for p := l; p != nil; p = p.Next {
		s = append(s, p.Val)
	}
	return s
}

func CompareInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}
