package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	cf "github.com/rishu/connect_four"
)

func main() {
	fmt.Println("=== Connect Four ===")
	fmt.Println()

	mode := selectGameMode()
	var p1, p2 *cf.Player
	var bot cf.BotEngine

	switch mode {
	case 1:
		p1 = cf.NewPlayer("Player 1", cf.Red)
		p2 = cf.NewPlayer("Player 2", cf.Yellow)
	case 2:
		p1 = cf.NewPlayer("Human", cf.Red)
		p2 = cf.NewPlayer("Bot (Simple)", cf.Yellow)
		bot = cf.NewSimpleBotEngine()
	case 3:
		p1 = cf.NewPlayer("Human", cf.Red)
		p2 = cf.NewPlayer("Bot (Smart)", cf.Yellow)
		bot = cf.NewSmartBotEngine()
	}

	game := cf.NewGame(p1, p2)
	playGame(game, bot)
}

func selectGameMode() int {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Select game mode:")
	fmt.Println("1. Human vs Human")
	fmt.Println("2. Human vs Simple Bot")
	fmt.Println("3. Human vs Smart Bot")
	fmt.Print("\nEnter choice (1-3): ")

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err == nil && choice >= 1 && choice <= 3 {
			return choice
		}
		fmt.Print("Invalid choice. Enter 1, 2, or 3: ")
	}
}

func playGame(game *cf.Game, bot cf.BotEngine) {
	reader := bufio.NewReader(os.Stdin)

	for !game.IsOver() {
		printBoard(game.Board())
		current := game.CurrentPlayer()
		fmt.Printf("\n%s's turn (%s)\n", current.Name(), current.Color())

		var col int
		if bot != nil && current == game.Player2() {
			col = bot.ChooseMove(game)
			fmt.Printf("Bot plays column %d\n", col)
		} else {
			col = getHumanMove(reader, game)
		}

		if col == -2 {
			if err := game.UndoLastMove(); err != nil {
				fmt.Printf("Cannot undo: %v\n", err)
			} else {
				fmt.Println("Move undone!")
			}
			continue
		}

		if err := game.MakeMoveByColumn(col); err != nil {
			fmt.Printf("Invalid move: %v\n", err)
		}
	}

	printBoard(game.Board())
	printResult(game)
}

func getHumanMove(reader *bufio.Reader, game *cf.Game) int {
	fmt.Printf("Enter column (0-6) or 'u' to undo: ")

	for {
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "u" || input == "undo" {
			return -2
		}

		col, err := strconv.Atoi(input)
		if err != nil {
			fmt.Print("Invalid input. Enter 0-6 or 'u': ")
			continue
		}

		if col < 0 || col > 6 {
			fmt.Print("Column must be 0-6: ")
			continue
		}

		if !game.Board().CanPlace(col) {
			fmt.Print("Column is full. Choose another: ")
			continue
		}

		return col
	}
}

func printBoard(board *cf.Board) {
	fmt.Println()
	fmt.Println(board.String())
}

func printResult(game *cf.Game) {
	fmt.Println()
	switch game.State() {
	case cf.Won:
		winner := game.Winner()
		fmt.Printf("Game Over! %s (%s) wins!\n", winner.Name(), winner.Color())
	case cf.Draw:
		fmt.Println("Game Over! It's a draw!")
	}
}
