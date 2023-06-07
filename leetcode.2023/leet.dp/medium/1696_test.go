package medium_test

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"testing"
)

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x int) {
	*h = append(*h, x)
}

func (h *IntHeap) Pop() int {
	top := (*h)[0]
	*h = (*h)[1:]
	return top
}

func maxResult(nums []int, k int) int {
	dp := make([]int, len(nums))

	for i := 0; i < len(nums); i++ {
		dp[i] = math.MinInt
	}

	dp[0] = nums[0]

	for i := 0; i < len(nums); i++ {
		for step := 1; step <= k; step++ {
			nextStep := i + step
			if nextStep >= len(nums) {
				break
			}

			sum := nums[nextStep] + dp[i]
			if sum > dp[nextStep] {
				dp[nextStep] = sum
			}
		}
	}

	return dp[len(nums)-1]
}

func Test_maxResult(t *testing.T) {
	type args struct {
		nums []int
		k    int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"negative", args{[]int{1, -1, -2, 4, -7, 3}, 2}, 7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := maxResult(tt.args.nums, tt.args.k); got != tt.want {
				t.Errorf("maxResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_largeSumResult(t *testing.T) {
	rawData, err := ioutil.ReadFile("1696_data.json")
	if err != nil {
		t.Fatal("Failed to load data file:", err)
	}

	var testData struct {
		Data []int `json:"data"`
		K    int   `json:"k"`
	}

	if err := json.Unmarshal(rawData, &testData); err != nil {
		t.Fatal("Failed to load json from test data:", err)
	}

	got := maxResult(testData.Data, testData.K)
	t.Log(got)
}
