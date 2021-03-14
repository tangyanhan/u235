package microsoft

import "testing"

func convert(s string, numRows int) string {
	if numRows == 1 {
		return s
	}
	// col * numRows + col * (numRows-2) > len(s)
	result := make([]byte, 0, len(s))
	for r := 0; r < numRows; r++ {
		step := (numRows - 1) * 2
		for j := 0; r+j*step-r*2 < len(s); j++ {
			colIndex := r + j*step
			if j != 0 && r > 0 && r+1 != numRows {
				result = append(result, s[colIndex-r*2])
			}
			if colIndex < len(s) {
				result = append(result, s[colIndex])
			}
		}
	}
	return string(result)
}

func Test_convert(t *testing.T) {
	type args struct {
		s       string
		numRows int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{
				s:       "PAYPALISHIRING",
				numRows: 3,
			},
			want: "PAHNAPLSIIGYIR",
		},
		{
			args: args{
				s:       "PAYPALISHIRING",
				numRows: 4,
			},
			want: "PINALSIGYAHRPI",
		},
		{
			args: args{
				s:       "A",
				numRows: 1,
			},
			want: "A",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convert(tt.args.s, tt.args.numRows); got != tt.want {
				t.Errorf("convert() = %v, want %v", got, tt.want)
			}
		})
	}
}
