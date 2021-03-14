package microsoft

const (
	Walked = byte(1)
	Sea    = '0'
	Land   = '1'
)

func isValid(grid [][]byte, i, j int) bool {
	return i >= 0 && j >= 0 && i < len(grid) && j < len(grid[0])
}

func walkPos(grid [][]byte, i, j int) {
	if !isValid(grid, i, j) {
		return
	}
	switch grid[i][j] {
	case Walked:
		return
	case Sea:
		return
	case Land:
		grid[i][j] = Walked
	}
	walkPos(grid, i+1, j)
	walkPos(grid, i, j+1)
	walkPos(grid, i, j-1)
	walkPos(grid, i-1, j)
}

func numIslands(grid [][]byte) int {
	var currentIsland int

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]); j++ {
			switch grid[i][j] {
			case Walked:
				continue
			case Sea:
				continue
			case Land:
				currentIsland++
				walkPos(grid, i, j)
			}
		}
	}
	return currentIsland
}
