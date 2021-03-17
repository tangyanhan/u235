package microsoft

func isValidSudoku(board [][]byte) bool {
	var rowSet [9]int
	var colSet [9]int
	var areaSet [3][3]int
	addNumToSet := func(set1, set2, set3 *int, num int) bool {
		mask := 1 << num
		if (mask&*set1 != 0) || (mask&*set2 != 0) || (mask&*set3 != 0) {
			return true
		}
		*set1 |= mask
		*set2 |= mask
		*set3 |= mask
		return false
	}
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[i]); j++ {
			if board[i][j] != '.' {
				v := int(board[i][j] - '0')
				if addNumToSet(&rowSet[i], &colSet[j], &areaSet[i/3][j/3], v) {
					return false
				}
			}
		}
	}
	return true
}
