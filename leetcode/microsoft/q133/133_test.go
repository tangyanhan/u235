package microsoft

type Node struct {
	Val       int
	Neighbors []*Node
}

func cloneGraph(node *Node) *Node {
	if node == nil {
		return nil
	}
	var nodeMap [101]*Node
	var visited [101]bool
	var getOrAddNode func(node *Node) *Node
	getOrAddNode = func(node *Node) *Node {
		clone := nodeMap[node.Val]
		if clone != nil {
			if visited[node.Val] {
				return clone
			}
			return clone
		}
		clone = &Node{
			Val:       node.Val,
			Neighbors: make([]*Node, len(node.Neighbors)),
		}
		nodeMap[node.Val] = clone
		visited[clone.Val] = true
		for i, neighbor := range node.Neighbors {
			clone.Neighbors[i] = getOrAddNode(neighbor)
		}
		return clone
	}
	clonedNode := getOrAddNode(node)

	return clonedNode
}
