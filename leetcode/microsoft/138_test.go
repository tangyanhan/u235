package microsoft

type RNode struct {
	Val    int
	Next   *RNode
	Random *RNode
}

func copyRandomList(head *RNode) *RNode {
	if head == nil {
		return nil
	}
	m := make(map[*RNode]*RNode)
	result := &RNode{}
	p := head
	cp := result
	m[p] = result
	for p != nil {
		cp.Val = p.Val
		m[p] = cp
		p = p.Next
		if p != nil {
			cp.Next = &RNode{}
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
