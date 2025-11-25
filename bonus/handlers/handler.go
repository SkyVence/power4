package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"power4/shared"
)

// GameData represents the data structure passed to the template
type GameData struct {
	Board         [][]int // Board (0=empty, 1=player1, 2=player2)
	CurrentPlayer int     // 1 or 2
	Player1Name   string  // Player 1 nickname
	Player2Name   string  // Player 2 nickname
	Player1Score  int     // Player 1's score
	Player2Score  int     // Player 2's score
	GameOver      bool    // Whether game is finished
	GameWon       bool    // Whether someone won
	GameDraw      bool    // Whether it's a draw
	Winner        int     // Winning player (1 or 2)
	ShowModal     bool    // Whether to show win/draw modal
	Message       string  // Status message to display
	ColumnIndices []int   // Column indices for iteration
	RowIndices    []int   // Row indices for iteration
	Rows          int     // Number of rows
	Columns       int     // Number of columns
	InverseGravity bool   // Whether gravity is currently inverted
	TurnCount     int     // Current turn count
}

// SetupData represents data for the setup page
type SetupData struct {
	Error string // Error message if any
}

// Extended game state with nicknames and custom features
type ExtendedGameState struct {
	game         *shared.Power
	player1Name  string
	player2Name  string
	player1Score int
	player2Score int
	turnCount    int // Track turns for gravity inversion
	inverseGravity bool // Current gravity state
}

// Global game state (in production, you'd use sessions or database)
var gameState *ExtendedGameState

func init() {
	// Initialize with default values (will be set from setup page)
	gameState = &ExtendedGameState{
		player1Name:  "Player 1",
		player2Name:  "Player 2",
		player1Score: 0,
		player2Score: 0,
		turnCount:    0,
		inverseGravity: false,
	}
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
	rows := gameState.game.Settings.Rows
	cols := gameState.game.Settings.Columns
	
	// Generate indices dynamically
	colIndices := make([]int, cols)
	rowIndices := make([]int, rows)
	for i := 0; i < cols; i++ {
		colIndices[i] = i
	}
	for i := 0; i < rows; i++ {
		rowIndices[i] = i
	}

	data := GameData{
		Board:          convertBoardToTemplate(gameState.game.Board),
		CurrentPlayer:  int(gameState.game.GetCurrentPlayer()) + 1, // Convert to 1-based
		Player1Name:    gameState.player1Name,
		Player2Name:    gameState.player2Name,
		Player1Score:   gameState.player1Score,
		Player2Score:   gameState.player2Score,
		GameOver:       gameState.game.IsGameOver(),
		GameWon:        false,
		GameDraw:       false,
		Winner:         0,
		ShowModal:      showModal,
		Message:        message,
		ColumnIndices:  colIndices,
		RowIndices:     rowIndices,
		Rows:           rows,
		Columns:        cols,
		InverseGravity: gameState.inverseGravity,
		TurnCount:      gameState.turnCount,
	}

	// Set game state specific fields
	if gameState.game.IsGameOver() {
		switch gameState.game.GetGameState() {
		case shared.BLUE_WINS:
			data.GameWon = true
			data.Winner = 1
			if showModal {
				gameState.player1Score++
			}
		case shared.RED_WINS:
			data.GameWon = true
			data.Winner = 2
			if showModal {
				gameState.player2Score++
			}
		case shared.DRAW:
			data.GameDraw = true
		}
		data.ShowModal = showModal
	}

	return data
}

// SetupHandler renders the setup page for nicknames and board size
func SetupHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("bonus/templates/setup.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := SetupData{}
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// StartGameHandler handles the game start with nicknames and board size
func StartGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	player1Name := r.FormValue("player1")
	player2Name := r.FormValue("player2")
	rowsStr := r.FormValue("rows")
	colsStr := r.FormValue("columns")

	// Validate inputs
	if player1Name == "" {
		player1Name = "Player 1"
	}
	if player2Name == "" {
		player2Name = "Player 2"
	}

	rows, err := strconv.Atoi(rowsStr)
	if err != nil || rows < 4 || rows > 15 {
		rows = 6 // Default
	}

	cols, err := strconv.Atoi(colsStr)
	if err != nil || cols < 4 || cols > 15 {
		cols = 7 // Default
	}

	// Set up game state
	gameState.player1Name = player1Name
	gameState.player2Name = player2Name
	gameState.player1Score = 0
	gameState.player2Score = 0
	gameState.turnCount = 0
	gameState.inverseGravity = false

	// Create new game with custom settings
	settings := shared.GameSettings{
		Rows:    rows,
		Columns: cols,
	}
	gameState.game = shared.NewGameInstance(settings)

	// Redirect to game page
	http.Redirect(w, r, "/bonus/game", http.StatusSeeOther)
}

// GameHandler renders the main game page
func GameHandler(w http.ResponseWriter, r *http.Request) {
	if gameState.game == nil {
		// No game initialized, redirect to setup
		http.Redirect(w, r, "/bonus/setup", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("bonus/templates/game.html")
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

// MakeMove handles column click moves with inverse gravity support
func MakeMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if gameState.game == nil {
		http.Redirect(w, r, "/bonus/setup", http.StatusSeeOther)
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

	// Check if move is valid (account for inverse gravity)
	var isValid bool
	if gameState.game.GetGameState() != shared.ONGOING {
		isValid = false
	} else if column < 0 || column >= gameState.game.Settings.Columns {
		isValid = false
	} else {
		// Check if ANY slot in the column is empty
		// A column is full only when ALL slots are filled
		hasEmptySlot := false
		for row := 0; row < gameState.game.Settings.Rows; row++ {
			if gameState.game.Board[row][column] == 0 {
				hasEmptySlot = true
				break
			}
		}
		isValid = hasEmptySlot
	}

	if !isValid {
		tmpl, err := template.ParseFiles("bonus/templates/game.html")
		if err != nil {
			http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var message string
		if gameState.game.IsGameOver() {
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

	// Make the move with inverse gravity support
	coord := shared.Coordinate{Column: column, Row: 0}
	makeMoveWithGravity(coord)

	// Increment turn count and check for gravity inversion (every 5 turns)
	gameState.turnCount++
	if gameState.turnCount%5 == 0 && !gameState.game.IsGameOver() {
		gameState.inverseGravity = !gameState.inverseGravity
	}

	// Check if this move ended the game
	showModal := gameState.game.IsGameOver()

	// Render the template with updated game state
	tmpl, err := template.ParseFiles("bonus/templates/game.html")
	if err != nil {
		http.Error(w, "Error parsing template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var message string
	if gameState.game.IsGameOver() {
		switch gameState.game.GetGameState() {
		case shared.BLUE_WINS:
			message = gameState.player1Name + " wins!"
		case shared.RED_WINS:
			message = gameState.player2Name + " wins!"
		case shared.DRAW:
			message = "It's a draw!"
		}
	} else if gameState.inverseGravity {
		message = "⚠️ Inverse Gravity Active! Pieces fall from bottom to top!"
	}

	data := createGameData(message, showModal)
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template: "+err.Error(), http.StatusInternalServerError)
	}
}

// makeMoveWithGravity makes a move considering inverse gravity
func makeMoveWithGravity(coord shared.Coordinate) {
	if gameState.inverseGravity {
		// Inverse gravity: pieces fall from bottom to top
		// This means pieces stack from top row downward (inverse of normal)
		// Find the lowest empty slot from top (first piece goes to top row)
		for row := 0; row < gameState.game.Settings.Rows; row++ {
			if gameState.game.Board[row][coord.Column] == 0 {
				// Place piece at the top-most empty position
				if gameState.game.IsPlaying == shared.BLUE {
					gameState.game.Board[row][coord.Column] = 'B'
				} else {
					gameState.game.Board[row][coord.Column] = 'R'
				}

				// Check for victory or draw after the move
				checkGameStateWithInverse(row, coord.Column)

				// Switch player only if game is still ongoing
				if gameState.game.GetGameState() == shared.ONGOING {
					if gameState.game.IsPlaying == shared.BLUE {
						gameState.game.IsPlaying = shared.RED
					} else {
						gameState.game.IsPlaying = shared.BLUE
					}
				}
				break
			}
		}
	} else {
		// Normal gravity (top to bottom) - pieces fall and stack from bottom
		gameState.game.MakeMove(coord)
	}
}

// checkGameStateWithInverse manually checks victory and draw for inverse gravity
func checkGameStateWithInverse(lastRow, lastCol int) {
	piece := gameState.game.Board[lastRow][lastCol]
	
	// Check for victory (4 in a row)
	if checkVictoryWithInverse(lastRow, lastCol, piece) {
		if gameState.game.IsPlaying == shared.BLUE {
			gameState.game.State = shared.BLUE_WINS
		} else {
			gameState.game.State = shared.RED_WINS
		}
		return
	}

	// Check for draw (board full)
	// In inverse gravity, pieces fill from top to bottom, so check bottom row
	// In normal gravity, pieces fill from bottom to top, so check top row
	checkRow := 0  // Default: normal gravity checks top row
	if gameState.inverseGravity {
		checkRow = gameState.game.Settings.Rows - 1  // Inverse gravity checks bottom row
	}
	
	isFull := true
	for col := 0; col < gameState.game.Settings.Columns; col++ {
		if gameState.game.Board[checkRow][col] == 0 {
			isFull = false
			break
		}
	}
	if isFull {
		gameState.game.State = shared.DRAW
	}
}

// checkVictoryWithInverse checks if there are 4 in a row
func checkVictoryWithInverse(row, col int, piece rune) bool {
	// Check horizontal
	if checkDirectionWithInverse(row, col, 0, 1, piece) >= 4 {
		return true
	}

	// Check vertical
	if checkDirectionWithInverse(row, col, 1, 0, piece) >= 4 {
		return true
	}

	// Check diagonal (top-left to bottom-right)
	if checkDirectionWithInverse(row, col, 1, 1, piece) >= 4 {
		return true
	}

	// Check diagonal (top-right to bottom-left)
	if checkDirectionWithInverse(row, col, 1, -1, piece) >= 4 {
		return true
	}

	return false
}

// checkDirectionWithInverse counts consecutive pieces in a given direction
func checkDirectionWithInverse(row, col, deltaRow, deltaCol int, piece rune) int {
	count := 1

	// Check in positive direction
	r, c := row+deltaRow, col+deltaCol
	for r >= 0 && r < gameState.game.Settings.Rows && c >= 0 && c < gameState.game.Settings.Columns && gameState.game.Board[r][c] == piece {
		count++
		r += deltaRow
		c += deltaCol
	}

	// Check in negative direction
	r, c = row-deltaRow, col-deltaCol
	for r >= 0 && r < gameState.game.Settings.Rows && c >= 0 && c < gameState.game.Settings.Columns && gameState.game.Board[r][c] == piece {
		count++
		r -= deltaRow
		c -= deltaCol
	}

	return count
}

// NewGameHandler starts a new game with same settings
func NewGameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if gameState.game == nil {
		http.Redirect(w, r, "/bonus/setup", http.StatusSeeOther)
		return
	}

	// Reset the game but keep nicknames and scores
	gameState.game.ResetGame()
	gameState.turnCount = 0
	gameState.inverseGravity = false

	// Redirect to game page
	http.Redirect(w, r, "/bonus/game", http.StatusSeeOther)
}

// ResetScoresHandler resets player scores
func ResetScoresHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	gameState.player1Score = 0
	gameState.player2Score = 0

	http.Redirect(w, r, "/bonus/game", http.StatusSeeOther)
}

