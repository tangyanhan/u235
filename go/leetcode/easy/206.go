package easy

func reverseList(head *ListNode) *ListNode {
	var n *ListNode
	for p := head; p != nil; {
		tmp := n
		n = p
		p = p.Next
		n.Next = tmp
	}
	return n
}
