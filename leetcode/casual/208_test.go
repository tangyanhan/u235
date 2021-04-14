package casual

import "testing"

type Trie struct {
	isEnd  bool
	childs [26]*Trie
}

/** Initialize your data structure here. */
func Constructor() Trie {
	return Trie{}
}

/** Inserts a word into the trie. */
func (t *Trie) Insert(word string) {
	curr := t
	for i := range word {
		c := word[i] - 'a'
		next := curr.childs[c]
		if next == nil {
			next = &Trie{}
			curr.childs[c] = next
		}
		curr = next
	}
	curr.isEnd = true
}

/** Returns if the word is in the trie. */
func (t *Trie) Search(word string) bool {
	curr := t
	for i := range word {
		c := word[i] - 'a'
		next := curr.childs[c]
		if next == nil {
			return false
		}
		curr = next
	}
	return curr.isEnd
}

/** Returns if there is any word in the trie that starts with the given prefix. */
func (t *Trie) StartsWith(prefix string) bool {
	curr := t
	for i := range prefix {
		c := prefix[i] - 'a'
		next := curr.childs[c]
		if next == nil {
			return false
		}
		curr = next
	}
	return true
}

func TestTrie(t *testing.T) {
	trie := Constructor()
	trie.Insert("apple")
	if !trie.Search("apple") {
		t.Fatal("no apple")
	}
	if trie.Search("app") {
		t.Fatal("should be no app")
	}
	if !trie.StartsWith("app") {
		t.Fatal("should have prefix app")
	}
}

/**
 * Your Trie object will be instantiated and called as such:
 * obj := Constructor();
 * obj.Insert(word);
 * param_2 := obj.Search(word);
 * param_3 := obj.StartsWith(prefix);
 */
