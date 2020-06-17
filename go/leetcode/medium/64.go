package medium

func minPathSum(grid [][]int) int {
	lastRow := len(grid) - 1
	lastCol := len(grid[0]) - 1
	for i := lastRow; i >= 0; i-- {
		for j := lastCol; j >= 0; j-- {
			if i == lastRow && j != lastCol {
				grid[i][j] = grid[i][j] + grid[i][j+1]
			} else if j == lastCol && i != lastRow {
				grid[i][j] = grid[i][j] + grid[i+1][j]
			} else if j != lastCol && i != lastRow {
				v1 := grid[i+1][j]
				v2 := grid[i][j+1]
				if v1 > v2 {
					grid[i][j] = grid[i][j] + v2
				} else {
					grid[i][j] = grid[i][j] + v1
				}
			}
		}
	}
	return grid[0][0]
}
