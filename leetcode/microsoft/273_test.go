package microsoft

import (
	"bytes"
	"testing"
)

var word0To9 = []string{
	"Zero",
	"One",
	"Two",
	"Three",
	"Four",
	"Five",
	"Six",
	"Seven",
	"Eight",
	"Nine",
}

var word10To19 = []string{
	"Ten",
	"Eleven",
	"Twelve",
	"Thirteen",
	"Fourteen",
	"Fifteen",
	"Sixteen",
	"Seventeen",
	"Eighteen",
	"Nineteen",
}

var wordTens = []string{
	"",
	"",
	"Twenty",
	"Thirty",
	"Forty",
	"Fifty",
	"Sixty",
	"Seventy",
	"Eighty",
	"Ninety",
}

var numSuffixes = []string{
	"",
	" Thousand",
	" Million",
	" Billion",
}

var strBuf [4096]byte

func numberToWords(num int) string {
	if num == 0 {
		return "Zero"
	}
	var parts []int
	for num != 0 {
		n := num % 1000
		num /= 1000
		parts = append(parts, n)
	}

	buf := bytes.NewBuffer(strBuf[:])
	buf.Reset()
	writeToBuf := func(s string) {
		if buf.Len() != 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(s)
	}
	for i := len(parts) - 1; i >= 0; i-- {
		below100 := parts[i] % 100
		hundred := parts[i] / 100
		if hundred != 0 {
			if buf.Len() != 0 {
				buf.WriteByte(' ')
			}
			buf.WriteString(word0To9[hundred])
			buf.WriteString(" Hundred")
		}
		if below100 > 0 {
			if below100 < 10 {
				writeToBuf(word0To9[below100])
			} else if below100 < 20 {
				writeToBuf(word10To19[below100-10])
			} else if below100 >= 20 {
				writeToBuf(wordTens[below100/10])
				digit := below100 % 10
				if digit != 0 {
					writeToBuf(word0To9[digit])
				}
			}
		}

		if parts[i] != 0 && numSuffixes[i] != "" {
			buf.WriteString(numSuffixes[i])
		}
	}
	return buf.String()
}

func Test_numberToWords(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want string
	}{
		{
			num:  1234567891,
			want: "One Billion Two Hundred Thirty Four Million Five Hundred Sixty Seven Thousand Eight Hundred Ninety One",
		},
		{
			num:  100,
			want: "One Hundred",
		},
		{
			num:  1000000,
			want: "One Million",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := numberToWords(tt.num); got != tt.want {
				t.Errorf("numberToWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
