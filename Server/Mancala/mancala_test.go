package Mancala

import (
	"fmt"
	"github.com/pkg/errors"
	"testing"
)

func failIfError(err error, t *testing.T) {
	if err != nil {
		err = errors.WithStack(err)
		t.Error(fmt.Sprintf("%+v", err))
		t.FailNow()
	}
}
func TestConstructor(t *testing.T) {
	board := NewBoard()
	if board == nil {
		t.Error("NewBoard has returned nil")
		t.FailNow()
	}
	if board.Result.GameOver {
		t.Error("Game is over before it even started")
	}
	for _, cell := range board.P1cells {
		if cell.Score != 4 {
			t.Error("Cell has wrong initial value")
		}
	}
	for _, cell := range board.P2cells {
		if cell.Score != 4 {
			t.Error("Cell has wrong initial value")
		}
	}
	if board.P1mc.Score != 0 || board.P2mc.Score != 0 {
		t.Error("Mancala has wrong initial value")
	}
}

func TestGame(t *testing.T) {
	board := NewBoard()
	board.Players[0] = "PlayerOne"
	board.Players[1] = "PlayerTwo"
	err := board.Turn(1, 2)
	failIfError(err, t)
	err = board.Turn(1, 5)
	failIfError(err, t)
	err = board.Turn(1, 0)
	if err == nil {
		t.Error("Expected error when making a turn out of order")
		t.Fail()
	}
	err = board.Turn(2, 3)
	failIfError(err, t)
	err = board.Turn(1, 3)
	failIfError(err, t)
	err = board.Turn(2, 1)
	failIfError(err, t)
	err = board.Turn(1, 1)
	failIfError(err, t)
	err = board.Turn(1, 0)
	failIfError(err, t)
	err = board.Turn(1, 4)
	failIfError(err, t)
	err = board.Turn(2, 0)
	failIfError(err, t)
	err = board.Turn(1, 5)
	failIfError(err, t)
	err = board.Turn(2, 0)
	failIfError(err, t)
	err = board.Turn(1, 0)
	failIfError(err, t)
	err = board.Turn(2, 2)
	failIfError(err, t)
	err = board.Turn(1, 3)
	failIfError(err, t)
	err = board.Turn(1, 4)
	failIfError(err, t)
	err = board.Turn(1, 5)
	failIfError(err, t)
	err = board.Turn(2, 1)
	failIfError(err, t)
	err = board.Turn(1, 0)
	failIfError(err, t)
	err = board.Turn(2, 0)
	failIfError(err, t)
	err = board.Turn(1, 2)
	failIfError(err, t)
	err = board.Turn(2, 5)
	failIfError(err, t)
	err = board.Turn(1, 1)
	failIfError(err, t)
	err = board.Turn(1, 3)
	failIfError(err, t)
	err = board.Turn(1, 0)
	failIfError(err, t)
	err = board.Turn(2, 0)
	failIfError(err, t)
	err = board.Turn(1, 4)
	failIfError(err, t)
	err = board.Turn(2, 2)
	failIfError(err, t)
	err = board.Turn(1, 5)
	failIfError(err, t)
	err = board.Turn(2, 2)
	failIfError(err, t)
	err = board.Turn(1, 2)
	failIfError(err, t)
	err = board.Turn(2, 0)
	failIfError(err, t)
	t.Log(board.String())
	t.Logf("checkIfOver returned %t", board.checkIfOver())
	if !board.checkIfOver() {
		fmt.Println(board.String())
		t.Error("Game should be ended by now")
		t.FailNow()
	}
}
