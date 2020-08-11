package main

import (
	"errors"
	"fmt"
	"strings"
)

type McCell struct {
	Score int `json:"Score"`
}

type GameResult struct {
	GameOver bool   `json:"game_over"`
	IsDraw   bool   `json:"is_draw"`
	Winner   string `json:"winner"`
}

type McBoard struct {
	P1score     int        `json:"player1_score"`
	P2score     int        `json:"player2_score"`
	P1cells     []McCell   `json:"player1_cells"`
	P2cells     []McCell   `json:"player2_cells"`
	P1mc        McCell     `json:"player1_mancala"`
	P2mc        McCell     `json:"player2_mancala"`
	NextPlayer  int        `json:"next_player"`
	Players     [2]string  `json:"players"`
	Result      GameResult `json:"result"`
	LastChanged int        `json:"last_changed"`
}

func (cell McCell) String() string {
	return fmt.Sprintf("(%2d)", cell.Score)
}
func (cell McCell) MancalaString() string {
	return fmt.Sprintf("[[%2d]]", cell.Score)
}
func (board McBoard) Printable(username string) string {
	length := len(board.P2cells)
	if length != len(board.P1cells) {
		panic(errors.New("invalid board state"))
	}
	var playerMC, opponentMC *McCell
	var playerCells, opponentCells []McCell
	if username == board.Players[0] {
		playerMC = &board.P1mc
		playerCells = board.P1cells

		opponentMC = &board.P2mc
		opponentCells = board.P2cells
	} else {
		playerMC = &board.P2mc
		playerCells = board.P2cells

		opponentMC = &board.P1mc
		opponentCells = board.P1cells
	}
	sb := strings.Builder{}
	sb.WriteString(opponentMC.MancalaString())
	for i := range board.P2cells {
		sb.WriteString(" ")
		sb.WriteString(opponentCells[length-1-i].String())
	}
	if board.NextPlayer > 0 {
		sb.WriteString("\t\tNext is:" + board.Players[board.NextPlayer-1])
	}
	sb.WriteString("\n")
	sb.WriteString("      ")
	for _, cell := range playerCells {
		sb.WriteString(" ")
		sb.WriteString(cell.String())
	}
	sb.WriteString(" ")
	sb.WriteString(playerMC.MancalaString())
	sb.WriteString("\n")
	return sb.String()
}
