package byted

import (
	"container/list"
	"testing"
)

type LRUCache struct {
	list     *list.List
	capacity int

	data [3000]*list.Element
}

type DataPair struct {
	Key   int16
	Value int16
}

func Constructor(capacity int) LRUCache {
	return LRUCache{
		capacity: capacity,
		list:     list.New(),
	}
}

func (this *LRUCache) Get(key int) int {
	node := this.data[key]
	if node == nil {
		return -1
	}
	this.list.MoveToFront(node)
	return int(node.Value.(*DataPair).Value)
}

func (this *LRUCache) Put(key int, value int) {
	node := this.data[key]
	if node == nil {
		node = this.list.PushFront(&DataPair{
			Key:   int16(key),
			Value: int16(value),
		})
		this.data[key] = node
		if this.list.Len() > this.capacity {
			back := this.list.Back()
			this.data[int(back.Value.(*DataPair).Key)] = nil
			this.list.Remove(back)
		}
		return
	}
	node.Value.(*DataPair).Value = int16(value)
	this.list.MoveToFront(node)
}

func TestConstructor(t *testing.T) {
	l := list.New()
	l.PushBack(13)
	type args struct {
		capacity int
		opts     [][]int
	}
	tests := []struct {
		name     string
		capacity int
		in       [][]int
		want     []int
	}{
		{
			capacity: 2,
			in:       [][]int{{1, 11}, {2, 22}, {1}, {3, 3}, {2}},
			want:     []int{11, -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lru := Constructor(tt.capacity)
			var idx int
			for _, v := range tt.in {
				if len(v) == 1 {
					value := lru.Get(v[0])
					if value != tt.want[idx] {
						t.Fatal("ops failed at index ", idx, " want:", tt.want[idx], "got:", value)
					}
					idx++
				} else {
					lru.Put(v[0], v[1])
				}
			}
		})
	}
}
