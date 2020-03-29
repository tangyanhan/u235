package easy

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	if l1 == nil {
		return l2
	}
	if l2 == nil {
		return l1
	}
	head, p := l1, l1
	if l1.Val > l2.Val {
		head, p = l2, l2
		l2 = l2.Next
	} else {
		l1 = l1.Next
	}
	for l1 != nil || l2 != nil {
		if l1 == nil {
			p.Next = l2
			return head
		}
		if l2 == nil {
			p.Next = l1
			return head
		}
		if l1.Val < l2.Val {
			p.Next = l1
			l1 = l1.Next
		} else {
			p.Next = l2
			l2 = l2.Next
		}
		p = p.Next
	}
	return head
}
