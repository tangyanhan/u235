package microsoft

import "testing"

type TicTacToe struct {
	data        [][]int
	playerCount [2]int
}

/** Initialize your data structure here. */
func Constructor(n int) TicTacToe {
	data := make([][]int, n)
	for i := 0; i < n; i++ {
		data[i] = make([]int, n)
	}
	return TicTacToe{
		data: data,
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
	this.data[row][col] = player
	this.playerCount[player-1]++
	cnt := this.playerCount[player-1]
	if cnt < len(this.data) {
		return 0
	}
	// scan row
	wins := true
	for i := 0; i < len(this.data); i++ {
		if this.data[row][i] != player {
			wins = false
			break
		}
	}
	if wins {
		return player
	}
	// scan columns
	wins = true
	for i := 0; i < len(this.data); i++ {
		if this.data[i][col] != player {
			wins = false
			break
		}
	}
	if wins {
		return player
	}
	// check diag
	if row+col+1 != len(this.data) && row != col {
		return 0
	}
	if row+col+1 == len(this.data) {
		wins = true
		for i := len(this.data) - 1; i >= 0; i-- {
			j := len(this.data) - 1 - i
			if this.data[i][j] != player {
				wins = false
				break
			}
		}
		if wins {
			return player
		}
	}
	if row == col {
		wins = true
		for i := 0; i < len(this.data); i++ {
			if this.data[i][i] != player {
				wins = false
				break
			}
		}
		if wins {
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
			for _, move := range tt.moves {
				if got := tt.this.Move(move.row, move.col, move.player); got != move.win {
					t.Errorf("TicTacToe.Move() = %v, want %v, data: %v", got, move.win, tt.this.data)
				}
			}

		})
	}
}
