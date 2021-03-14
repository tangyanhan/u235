package microsoft

type Node struct {
	Val    int
	Next   *Node
	Random *Node
}

func copyRandomList(head *Node) *Node {
	if head == nil {
		return nil
	}
	m := make(map[*Node]*Node)
	result := &Node{}
	p := head
	cp := result
	m[p] = result
	for p != nil {
		cp.Val = p.Val
		m[p] = cp
		p = p.Next
		if p != nil {
			cp.Next = &Node{}
			cp = cp.Next
		}
	}
	for p, cp := range m {
		if p.Random != nil {
			cp.Random = m[p.Random]
		}
	}
	return result
}
