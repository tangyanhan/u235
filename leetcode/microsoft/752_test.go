package microsoft

import (
	"container/list"
	"testing"
)

func openLockOptimized(deadends []string, target string) int {
	// 使用int代替string，使用Int16缩减内存占用
	toInt := func(s string) int16 {
		return int16(s[0]-'0')*1000 + int16(s[1]-'0')*100 + int16(s[2]-'0')*10 + int16(s[3]-'0')
	}
	spinNum := func(num int16) []int16 {
		d := [4]int16{
			num / 1000,
			(num % 1000) / 100,
			(num % 100) / 10,
			num % 10,
		}
		return []int16{
			(d[0]+1)%10*1000 + d[1]*100 + d[2]*10 + d[3],
			(d[0]-1+10)%10*1000 + d[1]*100 + d[2]*10 + d[3],
			d[0]*1000 + (d[1]+1)%10*100 + d[2]*10 + d[3],
			d[0]*1000 + (d[1]-1+10)%10*100 + d[2]*10 + d[3],
			d[0]*1000 + d[1]*100 + (d[2]+1)%10*10 + d[3],
			d[0]*1000 + d[1]*100 + (d[2]-1+10)%10*10 + +d[3],
			d[0]*1000 + d[1]*100 + d[2]*10 + (d[3]+1)%10,
			d[0]*1000 + d[1]*100 + d[2]*10 + (d[3]-1+10)%10,
		}
	}
	// 因为密码锁域是0000-9999，因此使用数组在时间空间上优于map
	var visited [10000]bool
	//visited := make(map[int16]bool)
	for _, s := range deadends {
		visited[toInt(s)] = true
	}

	targetInt := toInt(target)

	queue := list.New()
	queue.PushBack(int16(0))
	var depth int
	for queue.Len() != 0 {
		length := queue.Len()
		for i := 0; i < length; i++ {
			p := queue.Front()
			queue.Remove(p)
			value := p.Value.(int16)
			if visited[value] {
				continue
			}
			if value == targetInt {
				return depth
			}
			visited[value] = true
			nums := spinNum(value)
			for _, num := range nums {
				if !visited[num] {
					queue.PushBack(num)
				}
			}
		}
		depth++
	}
	return -1
}

func openLock(deadends []string, target string) int {
	visited := make(map[string]struct{})
	for _, s := range deadends {
		visited[s] = struct{}{}
	}

	queue := list.New()
	var depth int
	queue.PushBack("0000")
	for queue.Len() != 0 {
		length := queue.Len()
		for i := 0; i < length; i++ {
			p := queue.Front()
			queue.Remove(p)
			curr := p.Value.(string)
			if _, exists := visited[curr]; exists {
				continue
			}
			if curr == target {
				return depth
			}
			visited[curr] = struct{}{}
			for i := 0; i < 4; i++ {
				buf := []byte(curr)
				// up
				buf[i] = (curr[i]-'0'+1)%10 + '0'
				up := string(buf)

				if _, exists := visited[up]; !exists {
					queue.PushBack(up)
				}
				buf[i] = (curr[i]-'0'-1+10)%10 + '0'
				down := string(buf)
				if _, exists := visited[down]; !exists {
					queue.PushBack(down)
				}
			}
		}
		depth++
	}
	return -1
}

// A set containing int16, using struct{} saves 10% memory than bool
type Int16Set map[int16]struct{}

// Add element
func (s Int16Set) Add(v int16) {
	s[v] = struct{}{}
}

// Contains element
func (s Int16Set) Contains(v int16) bool {
	_, ok := s[v]
	return ok
}

func openLockDBFS(deadends []string, target string) int {
	// 使用int代替string，使用Int16缩减内存占用
	toInt := func(s string) int16 {
		return int16(s[0]-'0')*1000 + int16(s[1]-'0')*100 + int16(s[2]-'0')*10 + int16(s[3]-'0')
	}
	spinNum := func(num int16) []int16 {
		d := [4]int16{
			num / 1000,
			(num % 1000) / 100,
			(num % 100) / 10,
			num % 10,
		}
		return []int16{
			(d[0]+1)%10*1000 + d[1]*100 + d[2]*10 + d[3],
			(d[0]-1+10)%10*1000 + d[1]*100 + d[2]*10 + d[3],
			d[0]*1000 + (d[1]+1)%10*100 + d[2]*10 + d[3],
			d[0]*1000 + (d[1]-1+10)%10*100 + d[2]*10 + d[3],
			d[0]*1000 + d[1]*100 + (d[2]+1)%10*10 + d[3],
			d[0]*1000 + d[1]*100 + (d[2]-1+10)%10*10 + +d[3],
			d[0]*1000 + d[1]*100 + d[2]*10 + (d[3]+1)%10,
			d[0]*1000 + d[1]*100 + d[2]*10 + (d[3]-1+10)%10,
		}
	}

	// 因为密码锁域是0000-9999，因此使用数组在时间空间上优于map
	var visited [10000]bool
	//visited := make(map[int16]bool)
	for _, s := range deadends {
		visited[toInt(s)] = true
	}

	startSet := make(Int16Set)
	endSet := make(Int16Set)
	startSet.Add(int16(0))
	endSet.Add(toInt(target))
	var depth int
	for len(startSet) != 0 && len(endSet) != 0 {
		if len(startSet) > len(endSet) {
			tmp := startSet
			startSet = endSet
			endSet = tmp
		}
		toVisit := make(Int16Set)
		for value := range startSet {
			if visited[value] {
				continue
			}
			if endSet.Contains(value) {
				return depth
			}
			visited[value] = true
			nums := spinNum(value)
			for _, num := range nums {
				if !visited[num] {
					toVisit.Add(num)
				}
			}
		}
		depth++
		startSet = endSet
		endSet = toVisit
	}
	return -1
}

func Test_openLock(t *testing.T) {
	type args struct {
		deadends []string
		target   string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				deadends: []string{"8888"},
				target:   "0009",
			},
			want: 1,
		},
		{
			args: args{
				deadends: []string{"0201", "0101", "0102", "1212", "2002"},
				target:   "0202",
			},
			want: 6,
		},
		{
			args: args{
				deadends: []string{"0000"},
				target:   "8888",
			},
			want: -1,
		},
		{
			args: args{
				deadends: []string{"8887", "8889", "8878", "8898", "8788", "8988", "7888", "9888"},
				target:   "8888",
			},
			want: -1,
		},
		{
			args: args{
				deadends: []string{"1102", "1001", "0111", "0202", "1000"},
				target:   "0122",
			},
			want: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := openLockDBFS(tt.args.deadends, tt.args.target); got != tt.want {
				t.Errorf("openLock() = %v, want %v", got, tt.want)
			}
		})
	}
}
