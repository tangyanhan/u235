package microsoft

// UnionSet 用于并查集的合并与查找
type UnionSet struct {
	Data  []int
	Count int
}

// NewUnionSet 创建一个新的并集，初始集合为每个元素
func NewUnionSet(length int) *UnionSet {
	u := &UnionSet{
		Data:  make([]int, length),
		Count: length,
	}
	for i := range u.Data {
		u.Data[i] = i
	}
	return u
}

// Find parent of x
func (u *UnionSet) Find(x int) int {
	if u.Data[x] == x {
		return x
	}
	u.Data[x] = u.Find(u.Data[x])
	return u.Data[x]
}

// Join join another set
func (u *UnionSet) Join(x, y int) {
	px := u.Data[x]
	py := u.Data[y]

	// 已经是同一个集合，无需操作
	if px == py {
		return
	}
	u.Data[py] = px
	u.Count--
}
