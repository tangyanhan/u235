package microsoft

func maximalNetworkRank(n int, roads [][]int) int {
	g := make([][]bool, n)
	for i := range g {
		g[i] = make([]bool, n)
	}
	degree := make([]int, n)
	for _, p := range roads {
		a, b := p[0], p[1]
		g[a][b] = true
		g[b][a] = true
		degree[a]++
		degree[b]++
	}
	var maxRank int
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			v := degree[i] + degree[j]
			if g[i][j] {
				v--
			}
			if v > maxRank {
				maxRank = v
			}
		}
	}
	return maxRank
}
