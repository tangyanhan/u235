package microsoft

import (
	"testing"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// byteQueue 是用一个slice来存储byte队列的结构，避免使用list的存储和性能损耗
type byteQueue struct {
	data  []byte
	head  int
	tail  int
	count int
}

func (q *byteQueue) pushBack(v byte) {
	q.data[q.tail] = v
	q.count++
	q.tail = (q.tail + 1) % len(q.data)
}

func (q *byteQueue) pop() byte {
	v := q.data[q.head]
	q.head = (q.head + 1) % len(q.data)
	q.count--
	return v
}

func (q *byteQueue) size() int {
	return q.count
}

func alienOrder(words []string) string {
	if len(words) == 1 {
		return words[0]
	}
	// 字符范围只有小写英文字母
	var g [26][26]bool
	var inDegree [26]int
	vertexMap := make(map[byte]bool)
	for i := 1; i < len(words); i++ {
		prev := words[i-1]
		curr := words[i]
		minLen := min(len(prev), len(curr))
		var edgeAdded bool
		var j int
		for j = 0; j < minLen; j++ {
			from := prev[j] - 'a'
			to := curr[j] - 'a'
			// 不管怎样，先把点放进去，因为可能会出现["abc", "abd"]，这时ab其实是不确定的，但得放进去
			vertexMap[from] = true
			vertexMap[to] = true
			if curr[j] == prev[j] {
				continue
			}
			// 构图过程中出现了 a->b, b->a，直接结束
			if g[to][from] {
				return ""
			}
			// 避免加入重复边，导致入度计算错误
			if !g[from][to] {
				g[from][to] = true
				inDegree[to]++
			}
			edgeAdded = true
			break
		}
		// 前缀都一样，但前面的更长一些，则无效: ["abcd", "abc"]
		if !edgeAdded && len(prev) > len(curr) {
			return ""
		}
		// 把剩余多出来的字符放入图中（可能没有边） ["ab", "abcdefg"]
		for k := j; k < len(prev); k++ {
			vertexMap[prev[k]-'a'] = true
		}
		for k := j; k < len(curr); k++ {
			vertexMap[curr[k]-'a'] = true
		}
	}

	queue := &byteQueue{
		data: make([]byte, 26),
	}
	for i := range vertexMap {
		if inDegree[i] == 0 {
			queue.pushBack(i)
		}
	}
	result := make([]byte, 0, 26)
	for queue.size() != 0 {
		from := queue.pop()
		result = append(result, 'a'+byte(from))
		for i := range vertexMap {
			if g[from][i] {
				inDegree[i]--
				if inDegree[i] == 0 {
					queue.pushBack(i)
				}
			}
		}
	}
	if len(result) != len(vertexMap) {
		return ""
	}
	return string(result)
}

func Test_alienOrder(t *testing.T) {
	tests := []struct {
		name  string
		words []string
		want  string
	}{
		{
			words: []string{"wrt", "wrf", "er", "ett", "rftt"},
			want:  "wertf",
		},
		{
			words: []string{"z", "z"},
			want:  "z",
		},
		{
			words: []string{"z", "zabcdefg"},
			want:  "zabcdefg",
		},
		{
			words: []string{"ac", "ab", "zc", "zb"},
			want:  "acbz",
		},
		{
			words: []string{"wlnb"},
			want:  "wlnb",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alienOrder(tt.words); got != tt.want {
				t.Errorf("alienOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
