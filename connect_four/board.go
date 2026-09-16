package connectfour

import (
	"fmt"
	"strings"
)

const (
	DefaultRows = 6
	DefaultCols = 7
	WinLength   = 4
)

// Board encapsulates the game grid and all placement/win-detection logic.
// Game delegates grid operations here while handling turn management itself.
type Board struct {
	rows int
	cols int
	grid [][]DiscColor
}

// NewBoard creates a board with the specified dimensions.
// All cells start empty (zero value of DiscColor).
func NewBoard(rows, cols int) *Board {
	grid := make([][]DiscColor, rows)
	for i := range grid {
		grid[i] = make([]DiscColor, cols)
	}
	return &Board{rows: rows, cols: cols, grid: grid}
}

// NewDefaultBoard creates a standard 6x7 Connect Four board.
func NewDefaultBoard() *Board {
	return NewBoard(DefaultRows, DefaultCols)
}

func (b *Board) Rows() int {
	return b.rows
}

func (b *Board) Cols() int {
	return b.cols
}

// CanPlace checks if a disc can be placed in the given column.
// Returns false if column is out of bounds or already full.
func (b *Board) CanPlace(col int) bool {
	if col < 0 || col >= b.cols {
		return false
	}
	return b.grid[0][col] == Empty
}

// PlaceDisc drops a disc into the specified column.
// Returns the row where the disc landed, or -1 if placement failed.
func (b *Board) PlaceDisc(col int, color DiscColor) int {
	if !b.CanPlace(col) {
		return -1
	}

	// gravity: disc falls to the lowest empty row
	for row := b.rows - 1; row >= 0; row-- {
		if b.grid[row][col] == Empty {
			b.grid[row][col] = color
			return row
		}
	}
	return -1
}

// ClearCell removes a disc from the specified cell.
// Used by undo functionality.
func (b *Board) ClearCell(row, col int) bool {
	if !b.InBounds(row, col) {
		return false
	}
	b.grid[row][col] = Empty
	return true
}

// IsFull returns true if no more moves are possible.
func (b *Board) IsFull() bool {
	for col := 0; col < b.cols; col++ {
		if b.CanPlace(col) {
			return false
		}
	}
	return true
}

// CheckWin determines if placing a disc at (row, col) creates a winning line.
// Uses directional vectors to check horizontal, vertical, and both diagonals.
func (b *Board) CheckWin(row, col int, color DiscColor) bool {
	if !b.InBounds(row, col) {
		return false
	}
	if b.grid[row][col] != color {
		return false
	}

	// direction vectors: horizontal, vertical, diagonal-down, diagonal-up
	directions := [][2]int{
		{0, 1},  // horizontal (left-right)
		{1, 0},  // vertical (up-down)
		{1, 1},  // diagonal \
		{1, -1}, // diagonal /
	}

	for _, dir := range directions {
		count := 1 // start with the placed disc
		count += b.countInDirection(row, col, dir[0], dir[1], color)
		count += b.countInDirection(row, col, -dir[0], -dir[1], color)
		if count >= WinLength {
			return true
		}
	}
	return false
}

// countInDirection counts consecutive discs of the same color in one direction.
func (b *Board) countInDirection(row, col, dRow, dCol int, color DiscColor) int {
	count := 0
	r, c := row+dRow, col+dCol

	for b.InBounds(r, c) && b.grid[r][c] == color {
		count++
		r += dRow
		c += dCol
	}
	return count
}

// Cell returns the disc color at the specified position.
// Returns Empty if coordinates are out of bounds.
func (b *Board) Cell(row, col int) DiscColor {
	if !b.InBounds(row, col) {
		return Empty
	}
	return b.grid[row][col]
}

// InBounds checks if coordinates are within the board.
func (b *Board) InBounds(row, col int) bool {
	return row >= 0 && row < b.rows && col >= 0 && col < b.cols
}

// ValidMoves returns a slice of column indices where a disc can be placed.
func (b *Board) ValidMoves() []int {
	moves := make([]int, 0, b.cols)
	for col := 0; col < b.cols; col++ {
		if b.CanPlace(col) {
			moves = append(moves, col)
		}
	}
	return moves
}

// String returns a visual representation of the board.
func (b *Board) String() string {
	var sb strings.Builder

	// column numbers header
	for col := 0; col < b.cols; col++ {
		sb.WriteString(fmt.Sprintf(" %d", col))
	}
	sb.WriteString("\n")

	// grid rows (top to bottom)
	for row := 0; row < b.rows; row++ {
		for col := 0; col < b.cols; col++ {
			sb.WriteString(" ")
			sb.WriteString(b.grid[row][col].Symbol())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// Copy creates a deep copy of the board for simulation purposes.
func (b *Board) Copy() *Board {
	newBoard := NewBoard(b.rows, b.cols)
	for row := 0; row < b.rows; row++ {
		copy(newBoard.grid[row], b.grid[row])
	}
	return newBoard
}
