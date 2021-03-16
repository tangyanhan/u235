package microsoft

import (
	"testing"
)

type TicTacToe struct {
	rowSums      [][2]uint8
	colSums      [][2]uint8
	diagonalSums [2][2]uint8
}

/** Initialize your data structure here. */
func Constructor(n int) TicTacToe {
	return TicTacToe{
		rowSums: make([][2]uint8, n),
		colSums: make([][2]uint8, n),
	}
}

/** Player {player} makes a move at ({row}, {col}).
  @param row The row of the board.
  @param col The column of the board.
  @param player The player, can be either 1 or 2.
  @return The current winning condition, can be either:
          0: No one wins.
          1: Player 1 wins.
          2: Player 2 wins. */
func (this *TicTacToe) Move(row int, col int, player int) int {
	cntIdx := player - 1
	this.rowSums[row][cntIdx]++
	if int(this.rowSums[row][cntIdx]) == len(this.rowSums) {
		return player
	}
	this.colSums[col][cntIdx]++
	if int(this.colSums[col][cntIdx]) == len(this.colSums) {
		return player
	}
	if row == col {
		this.diagonalSums[0][cntIdx]++
		if int(this.diagonalSums[0][cntIdx]) == len(this.rowSums) {
			return player
		}
	}
	if row+col+1 == len(this.rowSums) {
		this.diagonalSums[1][cntIdx]++
		if int(this.diagonalSums[1][cntIdx]) == len(this.rowSums) {
			return player
		}
	}

	return 0
}

func TestTicTacToe_Move(t *testing.T) {
	type move struct {
		row    int
		col    int
		player int
		win    int
	}
	tests := []struct {
		name  string
		this  TicTacToe
		moves []move
		want  int
	}{
		{
			this: Constructor(3),
			moves: []move{
				{0, 0, 1, 0},
				{0, 2, 2, 0},
				{2, 2, 1, 0},
				{1, 1, 2, 0},
				{2, 0, 1, 0},
				{1, 0, 2, 0},
				{2, 1, 1, 1},
			},
		},
		{
			this: Constructor(2),
			moves: []move{
				{0, 0, 2, 0},
				{0, 1, 1, 0},
				{1, 1, 2, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, move := range tt.moves {
				if got := tt.this.Move(move.row, move.col, move.player); got != move.win {
					t.Errorf("TicTacToe.Move() = %v, want %v, move=#%d, move=%v", got, move.win, i, move)
				}
			}

		})
	}
}
