package medium

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
type ListNode struct {
	Val  int
	Next *ListNode
}

func CreateList(s []int) *ListNode {
	if len(s) == 0 {
		return nil
	}
	head := &ListNode{
		Val: s[0],
	}
	p := head
	for _, v := range s[1:] {
		p.Next = &ListNode{
			Val: v,
		}
		p = p.Next
	}
	return head
}

func ListToSlice(head *ListNode) []int {
	if head == nil {
		return nil
	}
	var s []int
	for head != nil {
		s = append(s, head.Val)
		head = head.Next
	}
	return s
}

func CmpSlice(a, b []int) bool {
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

func removeNthFromEnd(head *ListNode, n int) *ListNode {
	if head == nil {
		return nil
	}
	fast := head
	for i := 1; i < n; i++ {
		fast = fast.Next
	}
	slow := head
	var preSlow *ListNode
	for fast.Next != nil {
		fast = fast.Next
		preSlow = slow
		slow = slow.Next
	}
	if preSlow != nil {
		preSlow.Next = slow.Next
	} else {
		head = slow.Next
	}

	return head
}
