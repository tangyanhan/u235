package microsoft

// ListNode list
type ListNode struct {
	Val  int
	Next *ListNode
}

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}
	result := &ListNode{}
	var carry int
	p, q, r := l1, l2, result
	for p != nil || q != nil {
		if r.Next == nil {
			r.Next = &ListNode{}
			r = r.Next
		}
		if p != nil {
			r.Val += p.Val
			p = p.Next
		}
		if q != nil {
			r.Val += q.Val
			q = q.Next
		}
		r.Val += carry
		carry = r.Val / 10
		r.Val %= 10
	}
	if carry > 0 {
		r.Next = &ListNode{Val: carry}
	}
	return result.Next
}
