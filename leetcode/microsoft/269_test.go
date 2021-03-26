package microsoft

import (
	"container/list"
	"fmt"
	"testing"
)

func alienOrder(words []string) string {
	var g [26][26]bool
	var inDegree [26]int
	vertexMap := make(map[byte]bool)
	for i := 0; i < 100; i++ {
		var newEdgeAdded bool
		var lastByte byte
		// TODO: deal with remaining characters
		for _, word := range words {
			if i >= len(word) {
				continue
			}
			if lastByte == 0 {
				vertexMap[word[i]-'a'] = true
				lastByte = word[i]
				newEdgeAdded = true
				continue
			}
			fmt.Println("i=", i, "word=", word, string(lastByte), "->", string(word[i]))
			if word[i] != lastByte {
				fmt.Println("Add:", string(lastByte), "->", string(word[i]))
				from := lastByte - 'a'
				to := word[i] - 'a'
				if g[from][to] {
					continue
				}
				vertexMap[from] = true
				vertexMap[to] = true
				if g[to][from] {
					fmt.Println("Already got:", string(to+'a'), string(from+'a'))
					return ""
				}
				newEdgeAdded = true
				inDegree[to]++
				g[from][to] = true
				lastByte = word[i]
			}
		}
		if !newEdgeAdded {
			break
		}
	}
	queue := list.New()
	for i := range vertexMap {
		if inDegree[i] == 0 {
			queue.PushBack(i)
		}
	}
	result := make([]byte, 0, 26)
	for queue.Len() != 0 {
		from := queue.Front().Value.(byte)
		queue.Remove(queue.Front())
		result = append(result, 'a'+byte(from))
		for i := range vertexMap {
			if g[from][i] {
				inDegree[i]--
				if inDegree[i] == 0 {
					queue.PushBack(i)
				}
			}
		}
	}
	if len(result) != len(vertexMap) {
		fmt.Println("result=", string(result), "vertexLen=", len(vertexMap))
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := alienOrder(tt.words); got != tt.want {
				t.Errorf("alienOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}
