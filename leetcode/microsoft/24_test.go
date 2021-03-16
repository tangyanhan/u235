package microsoft

func swapPairs(head *ListNode) *ListNode {
	p := head
	for p != nil && p.Next != nil {
		p.Val, p.Next.Val = p.Next.Val, p.Val
		p = p.Next.Next
	}
	return head
}

func swapPairsRerverse(head *ListNode) *ListNode {
	p := head
	var prev *ListNode
	for p != nil && p.Next != nil {
		tmp := p.Next
		p.Next = p.Next.Next
		tmp.Next = p
		if prev != nil {
			prev.Next = tmp
		} else {
			head = tmp
		}

		prev = p
		p = p.Next
	}
	return head
}
