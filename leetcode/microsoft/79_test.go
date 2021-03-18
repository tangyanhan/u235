package microsoft

import (
	"testing"
)

func existBruteForce(board [][]byte, word string) bool {
	visited := make([]bool, len(board)*len(board[0]))
	var visit func(depth int, row, col int) bool
	visit = func(depth int, row, col int) bool {
		if row < 0 || col < 0 || row >= len(board) || col >= len(board[row]) {
			return false
		}

		if board[row][col] != word[depth] {
			return false
		}
		iv := row*len(board[0]) + col
		if visited[iv] {
			return false
		}
		if depth+1 == len(word) {
			return true
		}

		visited[iv] = true
		nextDepth := depth + 1
		if visit(nextDepth, row-1, col) ||
			visit(nextDepth, row, col+1) ||
			visit(nextDepth, row, col-1) ||
			visit(nextDepth, row+1, col) {
			return true
		}
		visited[iv] = false
		return false
	}
	for i := range board {
		for j := range board[i] {
			if board[i][j] == word[0] {
				if visit(0, i, j) {
					return true
				}
			}
		}
	}
	return false
}

func Test_exist(t *testing.T) {
	type args struct {
		board [][]byte
		word  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			args: args{
				board: [][]byte{{'A', 'B', 'C', 'E'}, {'S', 'F', 'C', 'S'}, {'A', 'D', 'E', 'E'}},
				word:  "ABCCED",
			},
			want: true,
		},
		{
			args: args{
				board: [][]byte{{'A', 'B', 'C', 'E'}, {'S', 'F', 'C', 'S'}, {'A', 'D', 'E', 'E'}},
				word:  "ABCB",
			},
			want: false,
		},
		{
			args: args{
				board: [][]byte{{'C', 'A', 'A'}, {'A', 'A', 'A'}, {'B', 'C', 'D'}},
				word:  "AAB",
			},
			want: true,
		},
		{
			args: args{
				// ABCE
				// SFES
				// ADEE
				// 0 1 2 3
				// 4 5 6 7
				// 8 9 10 11
				board: [][]byte{{'A', 'B', 'C', 'E'}, {'S', 'F', 'E', 'S'}, {'A', 'D', 'E', 'E'}},
				word:  "ABCESEEEFS",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := existBruteForce(tt.args.board, tt.args.word); got != tt.want {
				t.Errorf("exist() = %v, want %v", got, tt.want)
			}
		})
	}
}
