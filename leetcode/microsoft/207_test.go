package microsoft

import (
	"container/list"
	"testing"
)

type adjVertex struct {
	Val      int
	InDegree int
	Head     *adjLinkedList
	Tail     *adjLinkedList
}

type adjLinkedList struct {
	To   *adjVertex
	Next *adjLinkedList
}

func (v *adjVertex) addEdge(to *adjVertex) {
	node := &adjLinkedList{
		To: to,
	}
	if v.Head == nil {
		v.Head = node
		v.Tail = node
		return
	}
	v.Tail.Next = node
	v.Tail = node
}

func canFinish(numCourses int, prerequisites [][]int) bool {
	if len(prerequisites) == 0 {
		return true
	}
	// prerequisites <=5000, while numCourses <=10^5, V > E, better use adj
	// only create vertex for those in the prerequisites

	nodeMap := make(map[int]*adjVertex)
	getOrAddNode := func(val int) *adjVertex {
		node, ok := nodeMap[val]
		if ok {
			return node
		}
		node = &adjVertex{Val: val}
		nodeMap[val] = node
		return node
	}

	for _, pair := range prerequisites {
		from := getOrAddNode(pair[1])
		to := getOrAddNode(pair[0])
		to.InDegree++
		from.addEdge(to)
	}

	queue := list.New()
	// push vertex with in degree==0
	for _, v := range nodeMap {
		if v.InDegree == 0 {
			queue.PushBack(v)
		}
	}
	var doneNum int
	for queue.Len() != 0 {
		p := queue.Front().Value.(*adjVertex)
		queue.Remove(queue.Front())
		doneNum++
		for next := p.Head; next != nil; next = next.Next {
			next.To.InDegree--
			if next.To.InDegree == 0 {
				queue.PushBack(next.To)
			}
		}
	}
	return doneNum == len(nodeMap)
}

func Benchmark_canFinish(b *testing.B) {
	tests := []struct {
		name          string
		numCourses    int
		prerequisites [][]int
		want          bool
	}{
		{
			numCourses:    2,
			prerequisites: [][]int{{1, 0}},
			want:          true,
		},
		{
			numCourses:    2,
			prerequisites: [][]int{{1, 0}, {0, 1}},
			want:          false,
		},
		{
			numCourses:    3,
			prerequisites: [][]int{{1, 0}, {2, 1}, {0, 2}},
			want:          false,
		},
		{
			numCourses:    1,
			prerequisites: [][]int{},
			want:          true,
		},
		{
			numCourses:    5,
			prerequisites: [][]int{{1, 4}, {2, 4}, {3, 1}, {3, 2}},
			want:          true,
		},
	}
	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if got := canFinish(tt.numCourses, tt.prerequisites); got != tt.want {
					b.Errorf("canFinish() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
