package main

import "errors"

type McCell struct {
	Score    int `json:"Score"`
	owner    int8
	board    *McBoard
	next     *McCell
	opposite *McCell
}

type GameResult struct {
	isDraw bool
	winner int
}

type McBoard struct {
	P1score    int      `json:"player1_score"`
	P2score    int      `json:"player2_score"`
	P1cells    []McCell `json:"player1_cells"`
	P2cells    []McCell `json:"player2_cells"`
	P1mc       McCell   `json:"player1_mancala"`
	P2mc       McCell   `json:"player2_mancala"`
	NextPlayer int      `json:"next_player"`
}

func newBoard() McBoard {
	length := 6
	result := McBoard{
		P1score: 0,
		P2score: 0,
		P1cells: make([]McCell, length),
		P2cells: make([]McCell, length),
	}
	for i := 0; i < length-1; i++ {
		result.P1cells[i].next = &result.P1cells[i+1]
		result.P2cells[i].next = &result.P2cells[i+1]
		result.P1cells[i].opposite = &result.P2cells[length-i-1]
		result.P2cells[i].opposite = &result.P1cells[length-i-1]
		result.P1cells[i].Score = 4
		result.P2cells[i].Score = 4
		result.P1cells[i].owner = 1
		result.P2cells[i].owner = 2
	}
	result.P1cells[length-1].Score = 4
	result.P2cells[length-1].Score = 4
	result.P1cells[length-1].owner = 1
	result.P2cells[length-1].owner = 2
	result.P1cells[length-1].next = &result.P1mc
	result.P2cells[length-1].next = &result.P2mc
	result.P1cells[length-1].opposite = &result.P2cells[0]
	result.P2cells[length-1].opposite = &result.P2cells[0]
	result.P1mc.next = &result.P2cells[0]
	result.P2mc.next = &result.P1cells[0]
	return result
}

func (board McBoard) endGame() GameResult {
	board.P1score = board.P1mc.Score
	board.P2score = board.P2mc.Score
	for i := range board.P1cells {
		board.P1score += board.P1cells[i].Score
		board.P2score += board.P2cells[i].Score
	}
	winner := 0
	if board.P1score > board.P2score {
		winner = 1
	} else if board.P2score > board.P1score {
		winner = 2
	}
	return GameResult{
		isDraw: board.P1score == board.P2score,
		winner: winner,
	}
}

func (cell *McCell) move() bool {
	player := cell.owner
	left := cell.Score
	board := cell.board
	cell.Score = 0
	var current *McCell = cell
	for left > 0 {
		current = current.next
		current.Score++
		left--
	}
	if player == 1 && current == &board.P1mc {
		return true
	} else if player == 2 && current == &board.P2mc {
		return true
	}
	if current.owner == player && current.Score == 1 && current.opposite != nil {
		gained := current.Score + current.opposite.Score
		current.Score = 0
		current.opposite.Score = 0
		if player == 1 {
			board.P1mc.Score += gained
		} else {
			board.P2mc.Score += gained
		}
	}
	return false
}

func (board *McBoard) turn(player, cell int) error {
	if player != board.NextPlayer {
		return errors.New("another player's turn")
	}
	if player > 2 || player < 1 {
		return errors.New("invalid player id")
	}
	if cell >= len(board.P1cells) || cell < 0 {
		return errors.New("invalid cell id")
	}
	if player == 1 {
		if board.P1cells[cell].Score == 0 {
			return errors.New("selected an empty cell")
		}

		if !board.P1cells[cell].move() {
			board.NextPlayer = 2
		}
		return nil
	} else {
		if board.P2cells[cell].Score == 0 {
			return errors.New("selected an empty cell")
		}
		if !board.P2cells[cell].move() {
			board.NextPlayer = 1
		}
		return nil
	}
}
