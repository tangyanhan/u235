package sensetime

import (
	"strconv"
	"testing"
)

func isValidPart(s string) bool {
	if s[0] == '0' {
		return len(s) == 1
	}
	v, _ := strconv.Atoi(s)
	return v <= 255
}

func restoreIP(s string, remain int, base string, results *[]string) {
	if remain == 1 {
		if s != "" && isValidPart(s) {
			*results = append(*results, s+"."+base)
		}
		return
	}
	maxRemainLen := remain * 3
	if len(s) < remain || len(s) > maxRemainLen {
		return
	}
	for i := 1; i <= 3; i++ {
		if len(s) < i {
			break
		}
		p := s[len(s)-i:]
		if !isValidPart(p) {
			continue
		}
		if remain != 4 {
			restoreIP(s[:len(s)-i], remain-1, p+"."+base, results)
		} else {
			restoreIP(s[:len(s)-i], remain-1, p, results)
		}
	}
}

func restoreIpAddresses(s string) []string {
	if len(s) == 12 {
		p := []string{s[0:3], s[3:6], s[6:9], s[9:]}
		for _, v := range p {
			if !isValidPart(v) {
				return nil
			}
		}
		return []string{p[0] + "." + p[1] + "." + p[2] + "." + p[3]}
	}
	result := make([]string, 0)
	restoreIP(s, 4, "", &result)
	return result
}

func Test_restoreIpAddresses(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{
			in:   "25525511135",
			want: []string{"255.255.11.135", "255.255.111.35"},
		},
		{
			in:   "1111",
			want: []string{"1.1.1.1"},
		},
		{
			in:   "255255255255",
			want: []string{"255.255.255.255"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := restoreIpAddresses(tt.in)
			if len(got) != len(tt.want) {
				t.Errorf("restoreIpAddresses() = %v, want %v", got, tt.want)
			}
			m := make(map[string]bool)
			for _, v := range tt.want {
				m[v] = true
			}
			for _, v := range got {
				if !m[v] {
					t.Errorf("restoreIpAddresses() = %v, want %v, missing %v", got, tt.want, v)
				}
			}
		})
	}
}
