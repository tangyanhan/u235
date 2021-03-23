package microsoft

import (
	"container/list"
	"testing"
)

func ladderLengthSimple(beginWord string, endWord string, wordList []string) int {
	dict := make(map[string]struct{})
	for _, word := range wordList {
		dict[word] = struct{}{}
	}
	visited := make(map[string]struct{})
	mapContains := func(m map[string]struct{}, s string) bool {
		_, exists := m[s]
		return exists
	}
	isLadder := func(a, b string) bool {
		var hasDiff bool
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				if hasDiff {
					return false
				}
				hasDiff = true
			}
		}
		return hasDiff
	}
	queue := list.New()
	queue.PushBack(beginWord)
	length := 1
	for queue.Len() != 0 {
		queueLen := queue.Len()
		for i := 0; i < queueLen; i++ {
			p := queue.Front()
			queue.Remove(p)
			word := p.Value.(string)
			if mapContains(visited, word) {
				continue
			}
			if word == endWord {
				return length
			}
			visited[beginWord] = struct{}{}
			for dictWord := range dict {
				if isLadder(word, dictWord) {
					queue.PushBack(dictWord)
					delete(dict, dictWord)
				}
			}
		}
		length++
	}
	return 0
}

// StringSet set of string
type StringSet map[string]struct{}

// Add s
func (set StringSet) Add(s string) {
	set[s] = struct{}{}
}

// Contains s
func (set StringSet) Contains(s string) bool {
	_, ok := set[s]
	return ok
}

func ladderLengthDBFS(beginWord string, endWord string, wordList []string) int {
	dict := make(StringSet)
	for _, word := range wordList {
		dict.Add(word)
	}
	visited := make(StringSet)
	isLadder := func(a, b string) bool {
		var hasDiff bool
		for i := 0; i < len(a); i++ {
			if a[i] != b[i] {
				if hasDiff {
					return false
				}
				hasDiff = true
			}
		}
		return hasDiff
	}

	startSet := make(StringSet)
	endSet := make(StringSet)
	startSet.Add(beginWord)
	endSet.Add(endWord)
	if !dict.Contains(endWord) {
		return 0
	}
	length := 1
	for len(startSet) != 0 && len(endSet) != 0 {
		if len(startSet) > len(endSet) {
			tmp := startSet
			startSet = endSet
			endSet = tmp
		}
		toVisit := make(StringSet)
		for word := range startSet {
			if visited.Contains(word) {
				continue
			}
			if endSet.Contains(word) {
				return length
			}
			visited.Add(word)
			for dictWord := range dict {
				if !visited.Contains(dictWord) && isLadder(word, dictWord) {
					toVisit.Add(dictWord)
				}
			}
		}

		length++
		startSet = endSet
		endSet = toVisit
	}
	return 0
}

func Test_ladderLength(t *testing.T) {
	type args struct {
		beginWord string
		endWord   string
		wordList  []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				beginWord: "hit",
				endWord:   "cog",
				wordList:  []string{"hot", "dot", "dog", "lot", "log", "cog"},
			},
			want: 5,
		},
		{
			args: args{
				beginWord: "hit",
				endWord:   "cog",
				wordList:  []string{"hot", "dot", "dog", "lot", "log"},
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ladderLengthDBFS(tt.args.beginWord, tt.args.endWord, tt.args.wordList); got != tt.want {
				t.Errorf("ladderLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
