package dp

// https://leetcode-cn.com/problems/triangle/
func minimumTotal(triangle [][]int) int {
	for i := len(triangle) - 2; i >= 0; i-- {
		layer := triangle[i]
		btmLayer := triangle[i+1]
		for col := range layer {
			left := col
			right := col + 1
			if btmLayer[left] < btmLayer[right] {
				layer[col] += btmLayer[left]
			} else {
				layer[col] += btmLayer[right]
			}
		}
	}
	return triangle[0][0]
}
