package lru_test

import (
	"fmt"
	"reflect"
	"testing"
)

type Node struct {
	Prev *Node
	Next *Node
	K    int
	V    int

	Frequency int
}

type LinkedList struct {
	Head   *Node
	Tail   *Node
	length int
}

func NewLinkedList() *LinkedList {
	l := LinkedList{
		Head: &Node{},
		Tail: &Node{},
	}
	l.Head.Next = l.Tail
	l.Tail.Prev = l.Head
	return &l
}

func (l *LinkedList) Len() int {
	return l.length
}

func (l *LinkedList) PushFront(node *Node) {
	next := l.Head.Next
	l.Head.Next = node
	node.Prev = l.Head
	node.Next = next
	next.Prev = node
	l.length++
}

func (l *LinkedList) Remove(node *Node) {
	node.Prev.Next = node.Next
	node.Next.Prev = node.Prev
	l.length--
}

type LFUCache struct {
	remainCap    int
	minFrequency int
	freqMap      map[int]*LinkedList
	kvMap        map[int]*Node
}

func Constructor(capacity int) LFUCache {
	c := LFUCache{
		remainCap: capacity,
		freqMap:   make(map[int]*LinkedList),
		kvMap:     make(map[int]*Node),
	}
	c.freqMap[0] = NewLinkedList()
	return c
}

func (this *LFUCache) Get(key int) int {
	if node, ok := this.kvMap[key]; ok {
		this.UpdateNode(node)
		return node.V
	}
	return -1
}

func (this *LFUCache) UpdateNode(node *Node) {
	list := this.freqMap[node.Frequency]
	list.Remove(node)
	if list.Len() == 0 && node.Frequency == this.minFrequency {
		this.minFrequency = node.Frequency + 1
	}

	node.Frequency++
	newList, ok := this.freqMap[node.Frequency]
	if !ok {
		newList = NewLinkedList()
		this.freqMap[node.Frequency] = newList
	}
	newList.PushFront(node)
}

func (this *LFUCache) Put(key int, value int) {
	if node, ok := this.kvMap[key]; ok {
		this.UpdateNode(node)
		node.V = value
		return
	}

	if this.remainCap == 0 {
		list := this.freqMap[this.minFrequency]
		delete(this.kvMap, list.Tail.Prev.K)
		list.Remove(list.Tail.Prev)
	} else {
		this.remainCap--
	}

	node := &Node{K: key, V: value}
	this.minFrequency = 0
	this.freqMap[0].PushFront(node)
	this.kvMap[key] = node
}

func Test_LFUCache(t *testing.T) {
	//["LRUCache","put","put","put","put","get","get"]
	//[[2],[2,1],[1,1],[2,3],[4,1],[1],[2]]
	tests := []struct {
		ops    []string
		values [][]int
		expect []int
	}{
		{[]string{"LFUCache", "put", "put", "get", "put", "get", "get", "put", "get", "get", "get"}, [][]int{{2}, {1, 1}, {2, 2}, {1}, {3, 3}, {2}, {3}, {4, 4}, {1}, {3}, {4}}, []int{1, -1, 3, -1, 3, 4}},
	}

	for i, test := range tests {
		testName := fmt.Sprintf("Cache-%d", i)
		t.Run(testName, func(t *testing.T) {
			var cache LFUCache
			got := make([]int, 0)
			for x, op := range test.ops {
				switch op {
				case "LFUCache":
					cache = Constructor(test.values[x][0])
				case "put":
					cache.Put(test.values[x][0], test.values[x][1])
					t.Logf("Put key=%d, value=%d", test.values[x][0], test.values[x][1])
					//t.Log(cache.String())
				case "get":
					got = append(got, cache.Get(test.values[x][0]))
				}
			}

			if !reflect.DeepEqual(got, test.expect) {
				t.Fatalf("Expect:%v Got:%v", test.expect, got)
			}
		})
	}
}
