package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"power4/shared"
)

// GameData represents the data structure passed to the template
type GameData struct {
	Board         [][]int // 6x7 board (0=empty, 1=player1, 2=player2)
	CurrentPlayer int     // 1 or 2
	Player1Score  int     // Player 1's score
	Player2Score  int     // Player 2's score
	GameOver      bool    // Whether game is finished
	GameWon       bool    // Whether someone won
	GameDraw      bool    // Whether it's a draw
	Winner        int     // Winning player (1 or 2)
	ShowModal     bool    // Whether to show win/draw modal
	Message       string  // Status message to display
	ColumnIndices []int   // [0,1,2,3,4,5,6] for iteration
	RowIndices    []int   // [0,1,2,3,4,5] for iteration
}

// Global game state (in production, you'd use sessions or database)
var (
	game         *shared.Power
	player1Score int
	player2Score int
)

func init() {
	// Initialize game with standard Connect 4 dimensions
	settings := shared.GameSettings{
		Rows:    6,
		Columns: 7,
	}
	game = shared.NewGameInstance(settings)
}

// convertBoardToTemplate converts the game board to template-friendly format
func convertBoardToTemplate(gameBoard [][]rune) [][]int {
	board := make([][]int, len(gameBoard))
	for i := range gameBoard {
		board[i] = make([]int, len(gameBoard[i]))
		for j := range gameBoard[i] {
			switch gameBoard[i][j] {
			case 0: // Empty
				board[i][j] = 0
			case 'B': // Blue player (Player 1)
				board[i][j] = 1
			case 'R': // Red player (Player 2)
				board[i][j] = 2
			default:
				board[i][j] = 0
			}
		}
	}
	return board
}

// createGameData creates the GameData struct for template rendering
func createGameData(message string, showModal bool) GameData {
	data := GameData{
		Board:         convertBoardToTemplate(game.Board),
		CurrentPlayer: int(game.GetCurrentPlayer()) + 1, // Convert to 1-based
		Player1Score:  player1Score,
		Player2Score:  player2Score,
		GameOver:      game.IsGameOver(),
		GameWon:       false,
		GameDraw:      false,
		Winner:        0,
		ShowModal:     showModal,
		Message:       message,
		ColumnIndices: []int{0, 1, 2, 3, 4, 5, 6},
		RowIndices:    []int{0, 1, 2, 3, 4, 5},
	}

	// Set game state specific fields
	if game.IsGameOver() {
		switch game.GetGameState() {
		case shared.BLUE_WINS:
			data.GameWon = true
			data.Winner = 1
			if showModal {
				player1Score++
			}
		case shared.RED_WINS:
			data.GameWon = true
			data.Winner = 2
			if showModal {
				player2Score++
			}
		case shared.DRAW:
			data.GameDraw = true
		}
		data.ShowModal = showModal
	}

	return data
}

// HomeHandler renders the main game page
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("base/templates/index.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := createGameData("", false)

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// MoveHandler handles column click moves
func MoveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	columnStr := r.FormValue("column")
	if columnStr == "" {
		http.Error(w, "Column parameter required", http.StatusBadRequest)
		return
	}

	column, err := strconv.Atoi(columnStr)
	if err != nil {
		http.Error(w, "Invalid column number", http.StatusBadRequest)
		return
	}

	// Check if move is valid
	coord := shared.Coordinate{Column: column, Row: 0} // Row doesn't matter for this check
	if !game.IsValidMove(coord) {
		// Redirect back with error message
		tmpl, err := template.ParseFiles("base/templates/index.html")
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var message string
		if game.IsGameOver() {
			message = "Game is already over!"
		} else {
			message = "Column is full! Try another column."
		}

		data := createGameData(message, false)
		err = tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// Make the move
	game.MakeMove(coord)

	// Check if this move ended the game
	showModal := game.IsGameOver()

	// Render the template with updated game state
	tmpl, err := template.ParseFiles("base/templates/index.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var message string
	if game.IsGameOver() {
		switch game.GetGameState() {
		case shared.BLUE_WINS:
			message = "Player 1 (Blue) wins!"
		case shared.RED_WINS:
			message = "Player 2 (Red) wins!"
		case shared.DRAW:
			message = "It's a draw!"
		}
	}

	data := createGameData(message, showModal)
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

// NewGameHandler starts a new game
func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reset the game
	game.ResetGame()

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ResetScoresHandler resets player scores
func ResetScoresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Reset scores
	player1Score = 0
	player2Score = 0

	// Redirect to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
