package lru_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type Node struct {
	Prev *Node
	Next *Node
	K    int
	V    int
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

type LRUCache struct {
	capacity int
	kvArr    [10001]*Node
	list     *LinkedList
}

func Constructor(capacity int) LRUCache {
	c := LRUCache{
		capacity: capacity,
		list:     NewLinkedList(),
	}
	return c
}

func (l *LRUCache) Get(key int) int {
	if e := l.kvArr[key]; e != nil {
		l.list.Remove(e)
		l.list.PushFront(e)
		return e.V
	}

	return -1
}

func (l *LRUCache) Put(key int, value int) {
	if e := l.kvArr[key]; e != nil {
		l.list.Remove(e)
		l.list.PushFront(e)
		e.V = value
		return
	}

	if l.list.Len() == l.capacity {
		last := l.list.Tail.Prev
		l.list.Remove(last)
		l.kvArr[last.K] = nil
	}

	e := &Node{K: key, V: value}
	l.list.PushFront(e)
	l.kvArr[key] = e
}

func (l *LRUCache) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Capacity=%d\n", l.capacity))
	for k, v := range l.kvArr {
		if v == nil {
			continue
		}
		b.WriteString(fmt.Sprintln("Key=", k, "Value=", v.V))
	}
	for e := l.list.Head.Next; e != l.list.Tail; e = e.Next {
		b.WriteString(fmt.Sprintln("Node - key=", e.K, "value=", e.V))
	}
	return b.String()
}

func Test_LRUCache(t *testing.T) {
	//["LRUCache","put","put","put","put","get","get"]
	//[[2],[2,1],[1,1],[2,3],[4,1],[1],[2]]
	tests := []struct {
		ops    []string
		values [][]int
		expect []int
	}{
		{[]string{"LRUCache", "put", "put", "put", "put", "get", "get"}, [][]int{{2}, {2, 1}, {1, 1}, {2, 3}, {4, 1}, {1}, {2}}, []int{-1, 3}},
		{[]string{"LRUCache", "put", "get"}, [][]int{{1}, {2, 1}, {2}}, []int{1}},
		{[]string{"LRUCache", "put", "put", "put", "put", "put", "get", "put", "get", "get", "put", "get", "put", "put", "put", "get", "put", "get", "get", "get", "get", "put", "put", "get", "get", "get", "put", "put", "get", "put", "get", "put", "get", "get", "get", "put", "put", "put", "get", "put", "get", "get", "put", "put", "get", "put", "put", "put", "put", "get", "put", "put", "get", "put", "put", "get", "put", "put", "put", "put", "put", "get", "put", "put", "get", "put", "get", "get", "get", "put", "get", "get", "put", "put", "put", "put", "get", "put", "put", "put", "put", "get", "get", "get", "put", "put", "put", "get", "put", "put", "put", "get", "put", "put", "put", "get", "get", "get", "put", "put", "put", "put", "get", "put", "put", "put", "put", "put", "put", "put"},
			[][]int{{10}, {10, 13}, {3, 17}, {6, 11}, {10, 5}, {9, 10}, {13}, {2, 19}, {2}, {3}, {5, 25}, {8}, {9, 22}, {5, 5}, {1, 30}, {11}, {9, 12}, {7}, {5}, {8}, {9}, {4, 30}, {9, 3}, {9}, {10}, {10}, {6, 14}, {3, 1}, {3}, {10, 11}, {8}, {2, 14}, {1}, {5}, {4}, {11, 4}, {12, 24}, {5, 18}, {13}, {7, 23}, {8}, {12}, {3, 27}, {2, 12}, {5}, {2, 9}, {13, 4}, {8, 18}, {1, 7}, {6}, {9, 29}, {8, 21}, {5}, {6, 30}, {1, 12}, {10}, {4, 15}, {7, 22}, {11, 26}, {8, 17}, {9, 29}, {5}, {3, 4}, {11, 30}, {12}, {4, 29}, {3}, {9}, {6}, {3, 4}, {1}, {10}, {3, 29}, {10, 28}, {1, 20}, {11, 13}, {3}, {3, 12}, {3, 8}, {10, 9}, {3, 26}, {8}, {7}, {5}, {13, 17}, {2, 27}, {11, 15}, {12}, {9, 19}, {2, 15}, {3, 16}, {1}, {12, 17}, {9, 1}, {6, 19}, {4}, {5}, {5}, {8, 1}, {11, 7}, {5, 2}, {9, 28}, {1}, {2, 2}, {7, 4}, {4, 22}, {7, 24}, {9, 26}, {13, 28}, {11, 26}},
			[]int{-1, 19, 17, -1, -1, -1, 5, -1, 12, 3, 5, 5, 1, -1, 30, 5, 30, -1, -1, 24, 18, -1, 18, -1, 18, -1, 4, 29, 30, 12, -1, 29, 17, 22, 18, -1, 20, -1, 18, 18, 20}},
	}

	for i, test := range tests {
		testName := fmt.Sprintf("Cache-%d", i)
		t.Run(testName, func(t *testing.T) {
			var cache LRUCache
			got := make([]int, 0)
			for x, op := range test.ops {
				switch op {
				case "LRUCache":
					cache = Constructor(test.values[x][0])
				case "put":
					cache.Put(test.values[x][0], test.values[x][1])
					t.Logf("Put key=%d, value=%d", test.values[x][0], test.values[x][1])
					t.Log(cache.String())
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

/**
 * Your LRUCache object will be instantiated and called as such:
 * obj := Constructor(capacity);
 * param_1 := obj.Get(key);
 * obj.Put(key,value);
 */
