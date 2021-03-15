package microsoft

import (
	"sort"
	"testing"
)

// https://leetcode-cn.com/problems/meeting-rooms-ii/
// 符合贪心选择性质
func minMeetingRooms(intervals [][]int) int {
	sort.Slice(intervals, func(i, j int) bool {
		a, b := intervals[i], intervals[j]
		if a[0] != b[0] {
			return a[0] < b[0]
		}
		return a[1] < b[1]
	})
	rooms := make([][]int, 0)
	var roomCount int
	for _, vec := range intervals {
		var found bool
		for _, room := range rooms {
			if vec[0] >= room[1] {
				room[1] = vec[1]
				found = true
				break
			}
		}
		if !found {
			roomCount++
			rooms = append(rooms, vec)
		}
	}
	return roomCount
}

func minMeetingRoomsMemSaving(intervals [][]int) int {
	// 按开始时间排序，开始时间低的靠前
	// 若开始时间相同，按结束时间越早的靠前
	sort.Slice(intervals, func(i, j int) bool {
		a, b := intervals[i], intervals[j]
		if a[0] != b[0] {
			return a[0] < b[0]
		}
		return a[1] < b[1]
	})
	// 每当找到一个区间，就从前向后找已经分配的房间，由于已经是排好序的，第一个找到的一定是最优解
	// intervals传入后已经无用，可以复用用来存储room区间
	var roomCount int
	for _, vec := range intervals {
		var found bool
		for _, room := range intervals[:roomCount] {
			if vec[0] >= room[1] {
				room[1] = vec[1]
				found = true
				break
			}
		}
		if !found {
			intervals[roomCount] = vec
			roomCount++
		}
	}
	return roomCount
}

func Test_minMeetingRooms(t *testing.T) {
	tests := []struct {
		name string
		in   [][]int
		want int
	}{
		{
			in:   [][]int{{9, 10}, {4, 9}, {4, 17}},
			want: 2,
		},
		{
			in:   [][]int{{2, 11}, {6, 16}, {11, 16}},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := minMeetingRoomsMemSaving(tt.in); got != tt.want {
				t.Errorf("minMeetingRooms() = %v, want %v", got, tt.want)
			}
		})
	}
}
