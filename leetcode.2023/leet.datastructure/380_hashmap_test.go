package leetdatastructure_test

import "math/rand"

type RandomizedSet struct {
	data   []int
	idxMap map[int]int
}

func Constructor() RandomizedSet {
	return RandomizedSet{
		data:   make([]int, 0, 4096),
		idxMap: make(map[int]int),
	}
}

func (this *RandomizedSet) Insert(val int) bool {
	if idx, ok := this.idxMap[val]; ok && idx != -1 {
		return false
	}

	this.data = append(this.data, val)
	this.idxMap[val] = len(this.data) - 1

	return true
}

func (this *RandomizedSet) Remove(val int) bool {
	if idx, ok := this.idxMap[val]; ok && idx != -1 {
		this.idxMap[val] = -1
		lastVal := this.data[len(this.data)-1]
		if lastVal != val {
			this.data[idx] = lastVal
			this.idxMap[lastVal] = idx
		}

		this.data = this.data[:len(this.data)-1]
		return true
	}
	return false
}

func (this *RandomizedSet) GetRandom() int {
	idx := rand.Intn(len(this.data))
	return this.data[idx]
}
