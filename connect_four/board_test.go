package connectfour

import (
	"testing"
)

func TestNewBoard(t *testing.T) {
	board := NewBoard(6, 7)

	if board.Rows() != 6 {
		t.Errorf("expected 6 rows, got %d", board.Rows())
	}
	if board.Cols() != 7 {
		t.Errorf("expected 7 columns, got %d", board.Cols())
	}

	// all cells should be empty
	for row := 0; row < 6; row++ {
		for col := 0; col < 7; col++ {
			if board.Cell(row, col) != Empty {
				t.Errorf("cell (%d,%d) should be empty", row, col)
			}
		}
	}
}

func TestNewDefaultBoard(t *testing.T) {
	board := NewDefaultBoard()
	if board.Rows() != DefaultRows || board.Cols() != DefaultCols {
		t.Errorf("expected %dx%d board, got %dx%d", DefaultRows, DefaultCols, board.Rows(), board.Cols())
	}
}

func TestCanPlace(t *testing.T) {
	board := NewDefaultBoard()

	// all columns should accept placement initially
	for col := 0; col < board.Cols(); col++ {
		if !board.CanPlace(col) {
			t.Errorf("column %d should accept placement", col)
		}
	}

	// out of bounds
	if board.CanPlace(-1) {
		t.Error("column -1 should be invalid")
	}
	if board.CanPlace(7) {
		t.Error("column 7 should be invalid")
	}
}

func TestPlaceDisc(t *testing.T) {
	board := NewDefaultBoard()

	// disc should fall to bottom row
	row := board.PlaceDisc(0, Red)
	if row != 5 {
		t.Errorf("first disc should land at row 5, got %d", row)
	}
	if board.Cell(5, 0) != Red {
		t.Errorf("cell (5,0) should be Red")
	}

	// second disc in same column should stack
	row = board.PlaceDisc(0, Yellow)
	if row != 4 {
		t.Errorf("second disc should land at row 4, got %d", row)
	}
}

func TestPlaceDiscFullColumn(t *testing.T) {
	board := NewDefaultBoard()

	// fill column 0
	for i := 0; i < 6; i++ {
		row := board.PlaceDisc(0, Red)
		if row == -1 {
			t.Errorf("disc %d should have been placed", i)
		}
	}

	// column should now be full
	if board.CanPlace(0) {
		t.Error("column 0 should be full")
	}

	row := board.PlaceDisc(0, Yellow)
	if row != -1 {
		t.Error("placing in full column should return -1")
	}
}

func TestPlaceDiscInvalidColumn(t *testing.T) {
	board := NewDefaultBoard()

	if board.PlaceDisc(-1, Red) != -1 {
		t.Error("placing in column -1 should return -1")
	}
	if board.PlaceDisc(7, Red) != -1 {
		t.Error("placing in column 7 should return -1")
	}
}

func TestIsFull(t *testing.T) {
	board := NewDefaultBoard()

	if board.IsFull() {
		t.Error("empty board should not be full")
	}

	// fill the entire board
	for col := 0; col < board.Cols(); col++ {
		for row := 0; row < board.Rows(); row++ {
			board.PlaceDisc(col, Red)
		}
	}

	if !board.IsFull() {
		t.Error("completely filled board should be full")
	}
}

func TestClearCell(t *testing.T) {
	board := NewDefaultBoard()

	board.PlaceDisc(0, Red)
	if board.Cell(5, 0) != Red {
		t.Error("disc should be at (5,0)")
	}

	board.ClearCell(5, 0)
	if board.Cell(5, 0) != Empty {
		t.Error("cell should be empty after clear")
	}
}

func TestClearCellOutOfBounds(t *testing.T) {
	board := NewDefaultBoard()

	if board.ClearCell(-1, 0) {
		t.Error("clearing out of bounds should return false")
	}
	if board.ClearCell(0, -1) {
		t.Error("clearing out of bounds should return false")
	}
}

func TestCheckWinHorizontal(t *testing.T) {
	board := NewDefaultBoard()

	// place 4 in a row horizontally at bottom
	for col := 0; col < 4; col++ {
		board.PlaceDisc(col, Red)
	}

	// check win from any of the 4 positions
	if !board.CheckWin(5, 0, Red) {
		t.Error("should detect horizontal win from position (5,0)")
	}
	if !board.CheckWin(5, 3, Red) {
		t.Error("should detect horizontal win from position (5,3)")
	}
}

func TestCheckWinVertical(t *testing.T) {
	board := NewDefaultBoard()

	// place 4 in a column
	for i := 0; i < 4; i++ {
		board.PlaceDisc(0, Red)
	}

	if !board.CheckWin(2, 0, Red) {
		t.Error("should detect vertical win")
	}
}

func TestCheckWinDiagonalDown(t *testing.T) {
	board := NewDefaultBoard()

	// create diagonal \ pattern
	// column 0: 3 Yellow, then Red on top
	// column 1: 2 Yellow, then Red
	// column 2: 1 Yellow, then Red
	// column 3: Red

	board.PlaceDisc(0, Yellow)
	board.PlaceDisc(0, Yellow)
	board.PlaceDisc(0, Yellow)
	board.PlaceDisc(0, Red) // row 2, col 0

	board.PlaceDisc(1, Yellow)
	board.PlaceDisc(1, Yellow)
	board.PlaceDisc(1, Red) // row 3, col 1

	board.PlaceDisc(2, Yellow)
	board.PlaceDisc(2, Red) // row 4, col 2

	board.PlaceDisc(3, Red) // row 5, col 3

	if !board.CheckWin(5, 3, Red) {
		t.Error("should detect diagonal \\ win")
	}
}

func TestCheckWinDiagonalUp(t *testing.T) {
	board := NewDefaultBoard()

	// create diagonal / pattern
	board.PlaceDisc(3, Yellow)
	board.PlaceDisc(3, Yellow)
	board.PlaceDisc(3, Yellow)
	board.PlaceDisc(3, Red) // row 2, col 3

	board.PlaceDisc(2, Yellow)
	board.PlaceDisc(2, Yellow)
	board.PlaceDisc(2, Red) // row 3, col 2

	board.PlaceDisc(1, Yellow)
	board.PlaceDisc(1, Red) // row 4, col 1

	board.PlaceDisc(0, Red) // row 5, col 0

	if !board.CheckWin(5, 0, Red) {
		t.Error("should detect diagonal / win")
	}
}

func TestCheckWinNoWin(t *testing.T) {
	board := NewDefaultBoard()

	// place only 3 in a row
	for col := 0; col < 3; col++ {
		board.PlaceDisc(col, Red)
	}

	if board.CheckWin(5, 0, Red) {
		t.Error("3 in a row should not be a win")
	}
}

func TestCheckWinWrongColor(t *testing.T) {
	board := NewDefaultBoard()

	for col := 0; col < 4; col++ {
		board.PlaceDisc(col, Red)
	}

	// checking for Yellow win should return false
	if board.CheckWin(5, 0, Yellow) {
		t.Error("should not detect win for wrong color")
	}
}

func TestBoardValidMoves(t *testing.T) {
	board := NewDefaultBoard()

	moves := board.ValidMoves()
	if len(moves) != 7 {
		t.Errorf("expected 7 valid moves, got %d", len(moves))
	}

	// fill column 0
	for i := 0; i < 6; i++ {
		board.PlaceDisc(0, Red)
	}

	moves = board.ValidMoves()
	if len(moves) != 6 {
		t.Errorf("expected 6 valid moves after filling column 0, got %d", len(moves))
	}

	// column 0 should not be in valid moves
	for _, col := range moves {
		if col == 0 {
			t.Error("column 0 should not be a valid move")
		}
	}
}

func TestInBounds(t *testing.T) {
	board := NewDefaultBoard()

	if !board.InBounds(0, 0) {
		t.Error("(0,0) should be in bounds")
	}
	if !board.InBounds(5, 6) {
		t.Error("(5,6) should be in bounds")
	}
	if board.InBounds(-1, 0) {
		t.Error("(-1,0) should be out of bounds")
	}
	if board.InBounds(6, 0) {
		t.Error("(6,0) should be out of bounds")
	}
}

func TestBoardCopy(t *testing.T) {
	board := NewDefaultBoard()
	board.PlaceDisc(0, Red)
	board.PlaceDisc(1, Yellow)

	copy := board.Copy()

	// verify copy has same state
	if copy.Cell(5, 0) != Red {
		t.Error("copy should have Red at (5,0)")
	}

	// modify original, copy should be unchanged
	board.PlaceDisc(2, Red)
	if copy.Cell(5, 2) != Empty {
		t.Error("copy should not be affected by original")
	}
}

func TestBoardString(t *testing.T) {
	board := NewDefaultBoard()
	board.PlaceDisc(0, Red)
	board.PlaceDisc(0, Yellow)

	str := board.String()
	if str == "" {
		t.Error("board string should not be empty")
	}
}
