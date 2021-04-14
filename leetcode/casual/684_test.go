package casual

func findRedundantConnection(edges [][]int) []int {
	parent := make([]int, len(edges)+1)
	for i := range parent {
		parent[i] = i
	}
	find := func(p int) int {
		for p != parent[p] {
			parent[p] = parent[parent[p]]
			p = parent[p]
		}
		return p
	}
	union := func(x, y int) bool {
		px, py := find(x), find(y)
		if px == py {
			return true
		}
		parent[px] = py
		return false
	}

	for _, edge := range edges {
		if union(edge[0], edge[1]) {
			return edge
		}
	}
	return nil
}
