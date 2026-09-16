package connectfour

import "errors"

var (
	ErrGameOver       = errors.New("game is already over")
	ErrNotYourTurn    = errors.New("not your turn")
	ErrInvalidMove    = errors.New("invalid move")
	ErrNoMovesToUndo  = errors.New("no moves to undo")
)

// Game is the orchestration layer that manages turns, validates moves,
// and tracks game state. External code interacts with the game through this class.
type Game struct {
	board         *Board
	player1       *Player
	player2       *Player
	currentPlayer *Player
	state         GameState
	winner        *Player // nil until someone wins
	moveHistory   []Move
}

// NewGame creates a new game with two players.
// Player1 always moves first.
func NewGame(player1, player2 *Player) *Game {
	return &Game{
		board:         NewDefaultBoard(),
		player1:       player1,
		player2:       player2,
		currentPlayer: player1,
		state:         InProgress,
		winner:        nil,
		moveHistory:   make([]Move, 0),
	}
}

// NewGameWithBoard creates a game with a custom board (for testing or variants).
func NewGameWithBoard(player1, player2 *Player, board *Board) *Game {
	return &Game{
		board:         board,
		player1:       player1,
		player2:       player2,
		currentPlayer: player1,
		state:         InProgress,
		winner:        nil,
		moveHistory:   make([]Move, 0),
	}
}

// MakeMove attempts to place the current player's disc in the specified column.
// Returns an error if the move is invalid.
func (g *Game) MakeMove(player *Player, col int) error {
	if g.state != InProgress {
		return ErrGameOver
	}
	if player != g.currentPlayer {
		return ErrNotYourTurn
	}

	row := g.board.PlaceDisc(col, player.Color())
	if row == -1 {
		return ErrInvalidMove
	}

	g.moveHistory = append(g.moveHistory, NewMove(player, row, col))

	if g.board.CheckWin(row, col, player.Color()) {
		g.state = Won
		g.winner = player
	} else if g.board.IsFull() {
		g.state = Draw
	} else {
		g.switchTurn()
	}

	return nil
}

// MakeMoveByColumn is a convenience method that uses the current player.
func (g *Game) MakeMoveByColumn(col int) error {
	return g.MakeMove(g.currentPlayer, col)
}

// UndoLastMove reverts the most recent move.
// Restores the board, switches back to the previous player, and resets game state.
func (g *Game) UndoLastMove() error {
	if len(g.moveHistory) == 0 {
		return ErrNoMovesToUndo
	}

	lastMove := g.moveHistory[len(g.moveHistory)-1]
	g.moveHistory = g.moveHistory[:len(g.moveHistory)-1]

	g.board.ClearCell(lastMove.Row, lastMove.Col)

	g.currentPlayer = lastMove.Player

	// reset game state since we're undoing
	g.state = InProgress
	g.winner = nil

	return nil
}

func (g *Game) switchTurn() {
	if g.currentPlayer == g.player1 {
		g.currentPlayer = g.player2
	} else {
		g.currentPlayer = g.player1
	}
}

// Getters

func (g *Game) CurrentPlayer() *Player {
	return g.currentPlayer
}

func (g *Game) State() GameState {
	return g.state
}

func (g *Game) Winner() *Player {
	return g.winner
}

func (g *Game) Board() *Board {
	return g.board
}

func (g *Game) Player1() *Player {
	return g.player1
}

func (g *Game) Player2() *Player {
	return g.player2
}

func (g *Game) MoveCount() int {
	return len(g.moveHistory)
}

func (g *Game) LastMove() (Move, bool) {
	if len(g.moveHistory) == 0 {
		return Move{}, false
	}
	return g.moveHistory[len(g.moveHistory)-1], true
}

// IsOver returns true if the game has ended.
func (g *Game) IsOver() bool {
	return g.state.IsOver()
}

// ValidMoves returns columns where a move can be made.
func (g *Game) ValidMoves() []int {
	if g.IsOver() {
		return nil
	}
	return g.board.ValidMoves()
}

// Opponent returns the other player.
func (g *Game) Opponent(player *Player) *Player {
	if player == g.player1 {
		return g.player2
	}
	return g.player1
}
