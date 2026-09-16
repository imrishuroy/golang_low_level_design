# Connect Four -- Low Level Design Document

A comprehensive guide to Object-Oriented Design and Design Patterns using Connect Four as a case study.

## Table of Contents

1. [Problem Statement](#problem-statement)
2. [Requirements](#requirements)
3. [Core Entities](#core-entities)
4. [Class Design](#class-design)
5. [Design Patterns Used](#design-patterns-used)
   - [Orchestrator Pattern](#1-orchestrator-pattern)
   - [Facade Pattern](#2-facade-pattern)
   - [State Pattern](#3-state-pattern-via-enum)
   - [Strategy Pattern](#4-strategy-pattern)
   - [Value Object Pattern](#5-value-object-pattern)
   - [Single Responsibility Principle](#6-single-responsibility-principle-srp)
   - [Open/Closed Principle](#7-openclosed-principle-ocp)
   - [Data-Driven Algorithm](#8-data-driven-algorithm-direction-vectors)
   - [Command Pattern (Partial)](#9-command-pattern-partial---move-history)
   - [Factory Pattern](#10-factory-pattern-constructors)
   - [Null Object Avoidance](#11-null-object-avoidance-via-pointer-semantics)
6. [Implementation Details](#implementation-details)
7. [Anti-Patterns Avoided](#anti-patterns-avoided)
8. [Extensibility](#extensibility)
9. [Interview Tips](#interview-tips)

---

## Problem Statement

> "Build the object-oriented design for a two-player Connect Four game. Players take turns dropping discs into a 7-column, 6-row board. The first to align four of their own discs vertically, horizontally, or diagonally wins."

Connect Four is a two-player connection game where players alternate dropping colored discs into a vertical grid. The disc falls to the lowest available position in the chosen column. The objective is to connect four discs in a row -- horizontally, vertically, or diagonally.

---

## Requirements

### Functional Requirements

| # | Requirement |
|---|-------------|
| 1 | Two players take turns dropping discs into a 7-column, 6-row board |
| 2 | A disc falls to the lowest available row in the chosen column |
| 3 | Game ends when a player gets four discs in a row (any direction) -- they win |
| 4 | Game ends when the board is full with no winner -- it's a draw |
| 5 | Invalid moves must be rejected: full column, out of turn, after game over |

### Out of Scope

- UI/rendering (backend logic only)
- Concurrent games
- Move history visualization
- Configurable board size
- Network multiplayer

---

## Core Entities

### Entity Identification

From the requirements, we identify three core entities by looking for **nouns**:

| Entity | Responsibility |
|--------|---------------|
| **Game** | Orchestrator. Holds the Board, tracks turns, manages game state (in progress, won, draw), enforces rules. When a player moves, Game validates, tells Board to place, checks for win, switches turns. |
| **Board** | The 7x6 grid. Owns grid state, handles disc placement. Knows how to check if a column is full, where a disc falls, whether four are connected. Doesn't know whose turn it is. |
| **Player** | Represents a participant. Simple data holder with name and disc color. No game logic. |

### Entity Relationships

```
                    +------------------+
                    |      Game        |  <-- Orchestrator/Facade
                    |------------------|
                    | - board          |
                    | - player1        |
                    | - player2        |
                    | - currentPlayer  |
                    | - state          |
                    | - winner         |
                    | - moveHistory    |
                    +------------------+
                      /      |       \
                     v       v        v
              +-------+  +--------+  +--------+
              | Board |  | Player |  | Player |
              +-------+  +--------+  +--------+
                  |
                  v
           +-------------+
           | DiscColor[] |  (grid)
           +-------------+
```

---

## Class Design

### Complete Class Diagram

```
+------------------+       +------------------+       +------------------+
|    GameState     |       |    DiscColor     |       |      Move        |
|------------------|       |------------------|       |------------------|
| InProgress = 0   |       | Empty = 0        |       | + Player *Player |
| Won = 1          |       | Red = 1          |       | + Row int        |
| Draw = 2         |       | Yellow = 2       |       | + Col int        |
+------------------+       +------------------+       +------------------+

+------------------+       +------------------+
|     Player       |       |      Board       |
|------------------|       |------------------|
| - name string    |       | - rows int       |
| - color DiscColor|       | - cols int       |
|------------------|       | - grid [][]Color |
| + Name() string  |       |------------------|
| + Color() Color  |       | + CanPlace(col)  |
+------------------+       | + PlaceDisc(col) |
                           | + CheckWin(r,c)  |
                           | + IsFull()       |
                           | + ClearCell(r,c) |
                           | + Copy()         |
                           +------------------+

+----------------------+       +------------------+
|        Game          |       |   <<interface>>  |
|----------------------|       |    BotEngine     |
| - board *Board       |       |------------------|
| - player1 *Player    |       | + ChooseMove(g)  |
| - player2 *Player    |       +------------------+
| - currentPlayer      |              ^
| - state GameState    |              |
| - winner *Player     |    +---------+---------+
| - moveHistory []Move |    |         |         |
|----------------------|    v         v         v
| + MakeMove(p, col)   | +------+ +-------+ +--------+
| + UndoLastMove()     | |Simple| | Smart | | Random |
| + CurrentPlayer()    | | Bot  | |  Bot  | |  Bot   |
| + State()            | +------+ +-------+ +--------+
| + Winner()           |
| + ValidMoves()       |
+----------------------+
```

---

## Design Patterns Used

### 1. Orchestrator Pattern

**What it is:**
The Orchestrator Pattern centralizes the coordination logic between multiple components. One class (the orchestrator) manages the workflow and tells other components what to do and when. Components don't communicate directly with each other.

**Why we need it:**
Without an orchestrator, components would have tangled dependencies:
- Board would need to know about Players
- Players would need to know about Board
- State management would be scattered

**How Game acts as Orchestrator:**

```go
type Game struct {
    board         *Board    // Component 1: Grid management
    player1       *Player   // Component 2: Player identity
    player2       *Player   // Component 3: Player identity
    currentPlayer *Player   // Orchestrator tracks whose turn
    state         GameState // Orchestrator tracks game state
    winner        *Player   // Orchestrator tracks outcome
    moveHistory   []Move    // Orchestrator tracks history
}
```

**The Orchestration Flow:**

```go
func (g *Game) MakeMove(player *Player, col int) error {
    // Step 1: Validate game state (orchestrator's knowledge)
    if g.state != InProgress {
        return ErrGameOver
    }
    
    // Step 2: Validate turn order (orchestrator's knowledge)
    if player != g.currentPlayer {
        return ErrNotYourTurn
    }
    
    // Step 3: Delegate grid work to Board
    row := g.board.PlaceDisc(col, player.Color())
    if row == -1 {
        return ErrInvalidMove
    }
    
    // Step 4: Record history (orchestrator's responsibility)
    g.moveHistory = append(g.moveHistory, NewMove(player, row, col))
    
    // Step 5: Ask Board about win, update state
    if g.board.CheckWin(row, col, player.Color()) {
        g.state = Won
        g.winner = player
    } else if g.board.IsFull() {
        g.state = Draw
    } else {
        g.switchTurn()  // Step 6: Manage turn switching
    }
    
    return nil
}
```

**Visual Representation:**

```
Without Orchestrator:              With Orchestrator:
                                   
  Board <---> Player1                    Game (Orchestrator)
    ^           ^                       /    |    \
    |           |                      v     v     v
    v           v                   Board  Player1  Player2
  Player2 <---> ???                   |
                                      v
(Tangled dependencies)             GameState
                                   
                                   (Clean, centralized control)
```

**Key Characteristics:**

| Aspect | Implementation |
|--------|----------------|
| Centralized control | All game flow goes through Game |
| Components are independent | Board doesn't know about Players |
| Orchestrator holds state | Game tracks turns, winner, history |
| Orchestrator makes decisions | Game decides if move is valid |
| Components just do their job | Board just manages the grid |

**Real-World Analogy:**
Think of a restaurant:
- **Chef (Board)** -- Knows how to cook, doesn't know who ordered
- **Waiter (Player)** -- Has identity, takes orders
- **Manager (Game)** -- Orchestrates: validates orders, tells chef what to cook, tracks tables

---

### 2. Facade Pattern

**What it is:**
The Facade Pattern provides a simplified interface to a complex subsystem. External code interacts with the facade, not the internal components.

**Why we need it:**
Without a facade, external code would need to:
- Understand Board's grid indexing
- Manually track turns
- Coordinate multiple objects

**How Game acts as Facade:**

```go
// External code ONLY uses these simple methods:
game := NewGame(p1, p2)
err := game.MakeMove(p1, 3)
err = game.UndoLastMove()
state := game.State()
winner := game.Winner()
moves := game.ValidMoves()
```

**What external code does NOT do:**

```go
// WRONG: Bypassing the facade
game.board.PlaceDisc(3, Red)  // Skips validation!
game.board.ClearCell(5, 3)    // Breaks state consistency!
game.currentPlayer = p2       // Breaks turn order!
```

**Facade vs Orchestrator:**

| Pattern | Focus | Our Usage |
|---------|-------|-----------|
| Orchestrator | Internal coordination | Game coordinates Board, Players, State |
| Facade | External simplification | Game provides simple API to callers |

Game implements BOTH patterns -- it orchestrates internally and provides a facade externally.

---

### 3. State Pattern (via Enum)

**What it is:**
The State Pattern allows an object to alter its behavior when its internal state changes. We implement it using a type-safe enumeration instead of boolean flags.

**The Problem with Booleans:**

```go
// ANTI-PATTERN: Boolean flags
type Game struct {
    isOver    bool
    hasWinner bool
    isDraw    bool
    winner    *Player
}
```

This creates **8 possible combinations** (2^3), but only **3 are valid**:

| isOver | hasWinner | isDraw | Valid? | Meaning |
|--------|-----------|--------|--------|---------|
| false | false | false | Yes | In progress |
| true | true | false | Yes | Won |
| true | false | true | Yes | Draw |
| false | true | false | NO | Won but not over? |
| true | true | true | NO | Both won and draw? |
| ... | ... | ... | NO | 5 more invalid states |

**The Solution:**

```go
type GameState int

const (
    InProgress GameState = iota  // 0
    Won                          // 1
    Draw                         // 2
)
```

Now there are **exactly 3 states** -- impossible states are unrepresentable.

**State Transitions:**

```
                    +-------------+
                    | InProgress  |
                    +-------------+
                      /         \
           (4 in a row)         (board full)
                    /             \
                   v               v
              +-------+        +-------+
              |  Won  |        | Draw  |
              +-------+        +-------+
```

**Usage in Code:**

```go
func (g *Game) MakeMove(player *Player, col int) error {
    // Clean state check
    if g.state != InProgress {
        return ErrGameOver
    }
    
    // ... place disc ...
    
    // Clean state transitions
    if g.board.CheckWin(row, col, player.Color()) {
        g.state = Won
        g.winner = player
    } else if g.board.IsFull() {
        g.state = Draw
    }
}
```

**Key Principle:** *Make invalid states unrepresentable.*

---

### 4. Strategy Pattern

**What it is:**
The Strategy Pattern defines a family of algorithms, encapsulates each one, and makes them interchangeable. The algorithm can vary independently from clients that use it.

**Why we need it:**
We want multiple AI difficulty levels without:
- Changing the Game class for each new AI
- Hardcoding AI logic in one place
- Making it hard to test AIs

**The Interface:**

```go
type BotEngine interface {
    ChooseMove(game *Game) int
}
```

**The Implementations:**

```go
// Strategy 1: SimpleBotEngine
// Always picks the first valid column
type SimpleBotEngine struct{}

func (b *SimpleBotEngine) ChooseMove(game *Game) int {
    moves := game.ValidMoves()
    if len(moves) == 0 {
        return -1
    }
    return moves[0]
}

// Strategy 2: SmartBotEngine
// Win if possible, block opponent, prefer center
type SmartBotEngine struct{}

func (b *SmartBotEngine) ChooseMove(game *Game) int {
    // Priority 1: Win
    if col := b.findWinningMove(game, game.CurrentPlayer()); col != -1 {
        return col
    }
    // Priority 2: Block
    if col := b.findWinningMove(game, game.Opponent(...)); col != -1 {
        return col
    }
    // Priority 3: Center
    // Priority 4: Random
}

// Strategy 3: RandomBotEngine
// Random selection for testing
type RandomBotEngine struct{}

func (b *RandomBotEngine) ChooseMove(game *Game) int {
    moves := game.ValidMoves()
    return moves[rand.Intn(len(moves))]
}
```

**Usage (strategies are interchangeable):**

```go
var bot BotEngine

// Choose strategy at runtime
switch difficulty {
case "easy":
    bot = NewSimpleBotEngine()
case "medium":
    bot = NewSmartBotEngine()
case "hard":
    bot = NewMinimaxBotEngine()  // Future addition
}

// Usage is identical regardless of strategy
col := bot.ChooseMove(game)
game.MakeMoveByColumn(col)
```

**Class Diagram:**

```
       +------------------+
       |   <<interface>>  |
       |    BotEngine     |
       +------------------+
       | + ChooseMove(g)  |
       +------------------+
              ^
              |
    +---------+---------+
    |         |         |
    v         v         v
+--------+ +-------+ +--------+
| Simple | | Smart | | Random |
|  Bot   | |  Bot  | |  Bot   |
+--------+ +-------+ +--------+
```

**Benefits:**
- **Open/Closed**: Add MinimaxBot without changing existing code
- **Testable**: Test each bot in isolation
- **Swappable**: Change difficulty at runtime
- **Game unchanged**: Doesn't know which bot is used

---

### 5. Value Object Pattern

**What it is:**
Value Objects are small, immutable objects whose equality is based on their values, not their identity. They carry data but have no significant behavior.

**Our Value Objects:**

```go
// Player: Identity value object
type Player struct {
    name  string     // Unexported: immutable after creation
    color DiscColor  // Unexported: immutable after creation
}

func NewPlayer(name string, color DiscColor) *Player {
    return &Player{name: name, color: color}
}

// Only getters, no setters
func (p *Player) Name() string { return p.name }
func (p *Player) Color() DiscColor { return p.color }
```

```go
// Move: History record value object
type Move struct {
    Player *Player
    Row    int
    Col    int
}

func NewMove(player *Player, row, col int) Move {
    return Move{Player: player, Row: row, Col: col}
}
```

```go
// DiscColor: Enumeration value object
type DiscColor int

const (
    Empty  DiscColor = iota
    Red
    Yellow
)
```

**Characteristics:**

| Characteristic | Player | Move | DiscColor |
|----------------|--------|------|-----------|
| Immutable | Yes (unexported fields) | Yes (struct copy) | Yes (primitive) |
| No behavior | Yes (only getters) | Yes (only constructor) | Yes (only String) |
| Identity by value | Yes | Yes | Yes |

**Why Player is NOT an interface:**

```go
// ANTI-PATTERN: Over-abstraction
type Player interface {
    Name() string
    Color() DiscColor
    ChooseMove(game *Game) int  // WRONG: Mixing concerns!
}

type HumanPlayer struct { ... }
type BotPlayer struct { ... }
```

This is wrong because:
- A human player doesn't "choose" moves -- the UI does
- Player is just *identity* -- who they are, not what they do
- Decision-making belongs in `BotEngine`, not `Player`

**Correct separation:**
- `Player` = identity (who)
- `BotEngine` = decision-making (how)

---

### 6. Single Responsibility Principle (SRP)

**What it is:**
A class should have only one reason to change. Each class has one job.

**Our Application:**

| Class | Single Responsibility | Changes When |
|-------|----------------------|--------------|
| `Board` | Grid state and physics | Grid rules change |
| `Game` | Turn management and flow | Game rules change |
| `Player` | Player identity | Player representation changes |
| `BotEngine` | Move selection | AI strategy changes |
| `GameState` | State representation | State model changes |
| `Move` | Move recording | History format changes |

**Board's Single Responsibility:**

```go
// Board ONLY handles grid physics
func (b *Board) PlaceDisc(col int, color DiscColor) int {
    // Does: Check bounds, apply gravity, place disc
    // Does NOT: Check turns, update game state, record history
}

func (b *Board) CheckWin(row, col int, color DiscColor) bool {
    // Does: Scan for 4-in-a-row
    // Does NOT: Declare winner, end game
}
```

**Game's Single Responsibility:**

```go
// Game ONLY handles game flow
func (g *Game) MakeMove(player *Player, col int) error {
    // Does: Validate turn, coordinate components, update state
    // Does NOT: Know how gravity works, how win detection works
}
```

**Test for SRP Violation:**
> "If you have to use the word 'and' to describe what a class does, it probably does too much."

- Board: "Manages the grid" -- Good
- Game: "Manages turns and coordinates components" -- Acceptable (orchestrator role)
- BadClass: "Manages grid AND validates turns AND records history" -- Violation

---

### 7. Open/Closed Principle (OCP)

**What it is:**
Software entities should be open for extension but closed for modification. You can add new functionality without changing existing code.

**Our Application:**

**Adding a new bot (Open for Extension):**

```go
// Add MinimaxBotEngine without changing ANY existing code
type MinimaxBotEngine struct {
    depth int
}

func NewMinimaxBotEngine(depth int) *MinimaxBotEngine {
    return &MinimaxBotEngine{depth: depth}
}

func (b *MinimaxBotEngine) ChooseMove(game *Game) int {
    // Minimax algorithm implementation
    return b.minimax(game, b.depth, true)
}

// Game, Board, other bots -- UNCHANGED
```

**Why this works:**
- `BotEngine` interface defines the contract
- New bots implement the interface
- Game uses the interface, not concrete types

**Closed for Modification:**

```go
// Game code NEVER changes when adding bots
func playWithBot(game *Game, bot BotEngine) {
    col := bot.ChooseMove(game)  // Works with ANY bot
    game.MakeMoveByColumn(col)
}
```

---

### 8. Data-Driven Algorithm (Direction Vectors)

**What it is:**
Instead of writing separate code for each variant, parameterize a single algorithm with data. The algorithm stays the same; the data changes.

**The Anti-Pattern (Code Duplication):**

```go
// WRONG: Four separate methods doing the same thing
func (b *Board) checkHorizontal(row, col int, color DiscColor) bool {
    count := 1
    for c := col - 1; c >= 0 && b.grid[row][c] == color; c-- {
        count++
    }
    for c := col + 1; c < b.cols && b.grid[row][c] == color; c++ {
        count++
    }
    return count >= 4
}

func (b *Board) checkVertical(...) bool { /* Same logic, different direction */ }
func (b *Board) checkDiagonalDown(...) bool { /* Same logic, different direction */ }
func (b *Board) checkDiagonalUp(...) bool { /* Same logic, different direction */ }
```

**Problems:**
- 4x code duplication
- Fix a bug in one? Must fix in all four
- Easy to have inconsistencies

**The Pattern (Data-Driven):**

```go
func (b *Board) CheckWin(row, col int, color DiscColor) bool {
    // Directions as DATA, not separate methods
    directions := [][2]int{
        {0, 1},   // horizontal: row stays, col changes
        {1, 0},   // vertical: row changes, col stays  
        {1, 1},   // diagonal \: both increase
        {1, -1},  // diagonal /: row up, col down
    }

    // ONE algorithm handles all directions
    for _, dir := range directions {
        count := 1
        count += b.countInDirection(row, col, dir[0], dir[1], color)
        count += b.countInDirection(row, col, -dir[0], -dir[1], color)
        if count >= WinLength {
            return true
        }
    }
    return false
}

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
```

**Why this is NOT Strategy Pattern:**
- Strategy: Runtime-swappable behaviors
- Data-Driven: Same algorithm, parameterized by data
- Connect Four directions are FIXED -- they never change

**Benefits:**
- Fix bug once, fixed everywhere
- Add Connect-5 by changing `WinLength = 5`
- Clear, testable, maintainable
- No code duplication

---

### 9. Command Pattern (Partial) -- Move History

**What it is:**
The Command Pattern encapsulates a request as an object, allowing you to parameterize clients, queue requests, and support undoable operations.

**Our Partial Implementation:**

```go
// Move encapsulates a single action
type Move struct {
    Player *Player
    Row    int
    Col    int
}

// Game stores command history
type Game struct {
    moveHistory []Move  // Stack of commands
}
```

**Execute (MakeMove):**

```go
func (g *Game) MakeMove(player *Player, col int) error {
    row := g.board.PlaceDisc(col, player.Color())
    
    // Record the command
    g.moveHistory = append(g.moveHistory, NewMove(player, row, col))
    
    // ... check win/draw ...
}
```

**Undo (UndoLastMove):**

```go
func (g *Game) UndoLastMove() error {
    if len(g.moveHistory) == 0 {
        return ErrNoMovesToUndo
    }
    
    // Pop last command
    lastMove := g.moveHistory[len(g.moveHistory)-1]
    g.moveHistory = g.moveHistory[:len(g.moveHistory)-1]
    
    // Reverse the command's effect
    g.board.ClearCell(lastMove.Row, lastMove.Col)
    g.currentPlayer = lastMove.Player
    g.state = InProgress
    g.winner = nil
    
    return nil
}
```

**Why "Partial" Command Pattern:**
A full Command Pattern would have:

```go
type Command interface {
    Execute()
    Undo()
}

type PlaceDiscCommand struct {
    game   *Game
    player *Player
    col    int
    row    int  // Filled on execute
}

func (c *PlaceDiscCommand) Execute() { ... }
func (c *PlaceDiscCommand) Undo() { ... }
```

We simplified because:
- Only one type of command (place disc)
- Undo logic is simple enough inline
- No need for command queuing or redo

---

### 10. Factory Pattern (Constructors)

**What it is:**
The Factory Pattern provides an interface for creating objects without specifying their concrete classes. In Go, we use constructor functions.

**Our Factories:**

```go
// Simple Factory: NewPlayer
func NewPlayer(name string, color DiscColor) *Player {
    return &Player{name: name, color: color}
}

// Simple Factory: NewBoard (with defaults)
func NewBoard(rows, cols int) *Board {
    grid := make([][]DiscColor, rows)
    for i := range grid {
        grid[i] = make([]DiscColor, cols)
    }
    return &Board{rows: rows, cols: cols, grid: grid}
}

// Convenience Factory: NewDefaultBoard
func NewDefaultBoard() *Board {
    return NewBoard(DefaultRows, DefaultCols)  // 6x7
}

// Factory with dependencies: NewGame
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

// Factory for testing: NewGameWithBoard
func NewGameWithBoard(player1, player2 *Player, board *Board) *Game {
    return &Game{
        board:         board,  // Inject custom board
        player1:       player1,
        player2:       player2,
        currentPlayer: player1,
        state:         InProgress,
        winner:        nil,
        moveHistory:   make([]Move, 0),
    }
}
```

**Benefits:**
- Encapsulate object creation complexity
- Provide sensible defaults
- Allow dependency injection for testing

---

### 11. Null Object Avoidance (via Pointer Semantics)

**What it is:**
Instead of using special "null" values or sentinel objects, we use pointer semantics where `nil` has clear meaning.

**The Problem:**

```go
// ANTI-PATTERN: Empty struct as "no winner"
type Game struct {
    winner Player  // What represents "no winner"?
}

// Confusing: Is this the winner or not?
if game.winner == Player{} { ... }  // Empty struct comparison
```

**Our Solution:**

```go
type Game struct {
    winner *Player  // nil = no winner, non-nil = winner
}

// Clear semantics
if game.winner == nil {
    fmt.Println("No winner yet")
} else {
    fmt.Printf("%s wins!", game.winner.Name())
}
```

**Where we use this:**

| Field | Type | Nil Meaning |
|-------|------|-------------|
| `winner` | `*Player` | No winner (yet or ever) |
| `currentPlayer` | `*Player` | Always set (never nil) |
| Bot in main | `BotEngine` | Human vs Human mode |

---

## Implementation Details

### Win Detection Algorithm

```go
func (b *Board) CheckWin(row, col int, color DiscColor) bool {
    // Early validation
    if !b.InBounds(row, col) || b.grid[row][col] != color {
        return false
    }

    // Direction vectors: [dRow, dCol]
    directions := [][2]int{
        {0, 1},   // horizontal
        {1, 0},   // vertical
        {1, 1},   // diagonal \
        {1, -1},  // diagonal /
    }

    for _, dir := range directions {
        count := 1  // The placed disc
        count += b.countInDirection(row, col, dir[0], dir[1], color)
        count += b.countInDirection(row, col, -dir[0], -dir[1], color)
        if count >= WinLength {
            return true
        }
    }
    return false
}
```

### Disc Placement (Gravity)

```go
func (b *Board) PlaceDisc(col int, color DiscColor) int {
    if !b.CanPlace(col) {
        return -1
    }

    // Gravity: disc falls to lowest empty row
    for row := b.rows - 1; row >= 0; row-- {
        if b.grid[row][col] == Empty {
            b.grid[row][col] = color
            return row
        }
    }
    return -1
}
```

---

## Anti-Patterns Avoided

### 1. Boolean Flag State

```go
// AVOIDED
type Game struct {
    isOver, hasWinner, isDraw bool  // 8 combinations, 3 valid
}

// USED
type Game struct {
    state GameState  // 3 states, 3 valid
}
```

### 2. God Class

```go
// AVOIDED: Everything in one class
type Game struct {
    // Board logic
    // Player logic
    // Win detection
    // AI logic
    // History logic
    // ... 1000 lines
}

// USED: Separated responsibilities
type Game struct { ... }   // Orchestration
type Board struct { ... }  // Grid
type Player struct { ... } // Identity
type BotEngine interface { ... }  // AI
```

### 3. Over-Engineered Win Checking

```go
// AVOIDED
type WinChecker interface { Check(...) bool }
type HorizontalChecker struct{}
type VerticalChecker struct{}
// 4 classes for fixed geometry

// USED
directions := [][2]int{{0,1}, {1,0}, {1,1}, {1,-1}}
// Data, not classes
```

### 4. Player with Behavior

```go
// AVOIDED
type Player interface {
    ChooseMove(game *Game) int  // Mixing identity with behavior
}

// USED
type Player struct { name, color }  // Pure identity
type BotEngine interface { ChooseMove(game *Game) int }  // Pure behavior
```

---

## Extensibility

### Adding New AI (Strategy Pattern)

```go
// Just implement the interface
type MinimaxBot struct{ depth int }

func (b *MinimaxBot) ChooseMove(game *Game) int {
    return b.minimax(game, b.depth, true)
}

// No changes to Game, Board, or other bots
```

### Adding Undo (Command Pattern)

Already implemented via `moveHistory` and `UndoLastMove()`.

### Adding Replay

```go
func (g *Game) Replay() []Move {
    return g.moveHistory  // Already stored
}
```

### Configurable Board Size

Already supported:

```go
board := NewBoard(8, 8)  // 8x8 board
game := NewGameWithBoard(p1, p2, board)
```

---

## Interview Tips

### What's Expected at Each Level

| Level | Expectations |
|-------|-------------|
| **Junior** | Identify Board, Game, Player. Working placement and win detection. Handle basic edge cases. |
| **Mid** | Clean separation of concerns. Direction-vector win checking. Discuss extensibility without implementing. Mention 1-2 patterns. |
| **Senior** | Production-quality design. Proactively explain design decisions. Handle multiple extensibility questions. Name and explain patterns used. |

### Patterns to Mention by Level

| Level | Must Mention | Good to Mention |
|-------|--------------|-----------------|
| Junior | SRP | State enum vs booleans |
| Mid | SRP, Strategy, State | Orchestrator, Factory |
| Senior | All above + Facade, Command, OCP | Data-driven algorithm, composition over inheritance |

### Common Mistakes to Avoid

1. **Boolean flags for state** -- Use enums
2. **Over-engineering win checks** -- Direction vectors, not 4 classes
3. **God class** -- Don't put everything in Game
4. **Forgetting edge cases** -- Full column, game over, wrong turn
5. **Player with behavior** -- Keep identity separate from decision-making

### Questions to Ask Interviewer

1. How do players interact? (column number input)
2. What are all end conditions? (win + draw)
3. How to handle invalid moves? (return error)
4. Single game or multiple concurrent? (single)
5. Backend only or UI too? (backend)
6. Need undo or history? (clarify scope)

---

## Summary

### Patterns Used

| Pattern | Where | Purpose |
|---------|-------|---------|
| **Orchestrator** | Game class | Coordinates Board, Players, State |
| **Facade** | Game class | Simple external API |
| **State (Enum)** | GameState | Prevent invalid states |
| **Strategy** | BotEngine | Swappable AI algorithms |
| **Value Object** | Player, Move, DiscColor | Immutable data carriers |
| **SRP** | All classes | One responsibility each |
| **OCP** | BotEngine interface | Extend without modifying |
| **Data-Driven** | CheckWin directions | One algorithm, parameterized |
| **Command (Partial)** | Move + history | Enable undo |
| **Factory** | New* functions | Encapsulate creation |
| **Null Avoidance** | *Player pointers | Clear nil semantics |

### Key Design Principles

1. **Make invalid states unrepresentable** -- Enum over booleans
2. **Separate identity from behavior** -- Player vs BotEngine
3. **Orchestrate, don't tangle** -- Game coordinates components
4. **Parameterize, don't duplicate** -- Direction vectors
5. **Open for extension, closed for modification** -- Interface + implementations
