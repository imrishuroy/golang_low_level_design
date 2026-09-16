# Connect Four - Low Level Design in Go

A clean, idiomatic Go implementation of the classic Connect Four game, designed as a learning exercise for Low Level Design (LLD) interviews.

## Features

- **Clean Architecture**: Separation of concerns between Game (orchestration), Board (grid logic), and Player (data)
- **Win Detection**: Efficient directional-vector algorithm for horizontal, vertical, and diagonal wins
- **Undo Support**: Full move history with undo capability
- **Bot Opponents**: Pluggable AI with multiple difficulty levels
  - `SimpleBotEngine`: Picks first available column
  - `SmartBotEngine`: Attempts to win, blocks opponent, prefers center
  - `RandomBotEngine`: Random valid moves
- **Comprehensive Tests**: 51 tests covering all game logic

## Project Structure

```
connect_four/
├── cmd/
│   └── main.go          # Interactive demo application
├── docs/
│   └── DESIGN.md        # Detailed design documentation
├── board.go             # Board state and win detection
├── board_test.go
├── bot.go               # AI opponent implementations
├── bot_test.go
├── disc_color.go        # Disc color enum
├── game.go              # Game orchestration and rules
├── game_test.go
├── game_state.go        # Game state enum
├── move.go              # Move record for history
├── player.go            # Player representation
├── go.mod
└── README.md
```

## Installation

```bash
go get github.com/rishu/connect_four
```

## Usage

### As a Library

```go
package main

import (
    cf "github.com/rishu/connect_four"
    "fmt"
)

func main() {
    // Create players
    p1 := cf.NewPlayer("Alice", cf.Red)
    p2 := cf.NewPlayer("Bob", cf.Yellow)

    // Start a new game
    game := cf.NewGame(p1, p2)

    // Make moves
    game.MakeMove(p1, 3)  // Player 1 drops in column 3
    game.MakeMove(p2, 4)  // Player 2 drops in column 4

    // Check game state
    fmt.Println(game.Board())
    fmt.Printf("Current player: %s\n", game.CurrentPlayer().Name())

    // Undo last move
    game.UndoLastMove()
}
```

### With Bot Opponent

```go
// Create a smart bot
bot := cf.NewSmartBotEngine()

// Get bot's move
col := bot.ChooseMove(game)
game.MakeMoveByColumn(col)
```

### Running the Demo

```bash
go run ./cmd/

# Or build and run
go build -o connect_four ./cmd/
./connect_four
```

## API Reference

### Game

| Method | Description |
|--------|-------------|
| `NewGame(p1, p2 *Player) *Game` | Create a new game |
| `MakeMove(player *Player, col int) error` | Make a move |
| `MakeMoveByColumn(col int) error` | Make a move using current player |
| `UndoLastMove() error` | Undo the last move |
| `CurrentPlayer() *Player` | Get current player |
| `State() GameState` | Get game state (InProgress, Won, Draw) |
| `Winner() *Player` | Get winner (nil if none) |
| `IsOver() bool` | Check if game has ended |
| `ValidMoves() []int` | Get available columns |

### Board

| Method | Description |
|--------|-------------|
| `NewBoard(rows, cols int) *Board` | Create custom board |
| `NewDefaultBoard() *Board` | Create standard 6x7 board |
| `CanPlace(col int) bool` | Check if column has space |
| `PlaceDisc(col int, color DiscColor) int` | Place disc, returns row |
| `CheckWin(row, col int, color DiscColor) bool` | Check for win |
| `Cell(row, col int) DiscColor` | Get cell contents |
| `Copy() *Board` | Deep copy for simulation |

### BotEngine Interface

```go
type BotEngine interface {
    ChooseMove(game *Game) int
}
```

Implementations:
- `SimpleBotEngine` - First valid column
- `SmartBotEngine` - Win/block/center strategy
- `RandomBotEngine` - Random selection

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run with coverage
go test -cover ./...
```

## Design Patterns Used

1. **State Pattern**: `GameState` enum cleanly represents game states
2. **Strategy Pattern**: `BotEngine` interface for swappable AI
3. **Value Objects**: `Move`, `Player` as immutable data holders
4. **Separation of Concerns**: Board handles grid, Game handles rules

## Key Design Decisions

1. **Enum over Booleans**: `GameState` prevents invalid state combinations
2. **Pointer for Optional**: `winner *Player` is nil for no winner (not empty struct)
3. **Direction Vectors**: Single algorithm handles all 4 win directions
4. **Board Owns Win Logic**: Game asks Board, doesn't know grid internals

## License

MIT
