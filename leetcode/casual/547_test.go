package casual

func findCircleNum(isConnected [][]int) int {
	parent := make([]int, len(isConnected))
	for i := range parent {
		parent[i] = i
	}
	cnt := len(parent)
	find := func(p int) int {
		for p != parent[p] {
			parent[p] = parent[parent[p]]
			p = parent[p]
		}
		return p
	}
	union := func(x, y int) {
		px, py := find(x), find(y)
		if px != py {
			parent[px] = py
			cnt--
		}
	}
	for i, row := range isConnected {
		for j, v := range row {
			if v == 1 {
				union(i, j)
			}
		}
	}
	return cnt
}
