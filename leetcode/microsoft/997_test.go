package microsoft

func findJudge(N int, trust [][]int) int {
	inDegree := make([]int, N+1)
	outDegree := make([]int, N+1)
	for _, pair := range trust {
		inDegree[pair[1]]++
		outDegree[pair[0]]++
	}
	for i := 1; i <= N; i++ {
		if inDegree[i] == N-1 && outDegree[i] == 0 {
			return i
		}
	}
	return -1
}
