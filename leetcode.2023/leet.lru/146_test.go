package lru_test

import (
	"container/list"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type LRUCache struct {
	capacity int
	kvArr    [10001]*list.Element
	kvMap    map[int]*list.Element
	list     list.List
}

type KV struct {
	key   int
	value int
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		kvMap:    make(map[int]*list.Element),
	}
}

func (l *LRUCache) Get(key int) int {
	if e, ok := l.kvMap[key]; ok && e != nil {
		l.list.MoveToFront(e)
		kv, _ := e.Value.(*KV)
		return kv.value
	}

	return -1
}

func (l *LRUCache) Put(key int, value int) {
	if e, ok := l.kvMap[key]; ok && e != nil {
		l.list.MoveToFront(e)
		kv, _ := e.Value.(*KV)
		kv.value = value
		return
	}

	if l.capacity <= 0 {
		if l.list.Back() != nil {
			kv := l.list.Back().Value.(*KV)
			l.kvMap[kv.key] = nil
			l.list.Remove(l.list.Back())
		}
	} else {
		l.capacity--
	}

	e := l.list.PushFront(&KV{key, value})
	l.kvMap[key] = e
}

func (l *LRUCache) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Capacity=%d\n", l.capacity))
	for k, v := range l.kvMap {
		if v == nil {
			continue
		}
		b.WriteString(fmt.Sprintln("Key=", k, "Value=", v.Value.(*KV).value))
	}
	for e := l.list.Front(); e != nil; e = e.Next() {
		kv := e.Value.(*KV)
		b.WriteString(fmt.Sprintln("Node - key=", kv.key, "value=", kv.value))
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
