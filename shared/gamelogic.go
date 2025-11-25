package shared

type Player int

const (
	BLUE Player = iota
	RED
)

type GameState int

const (
	ONGOING GameState = iota
	BLUE_WINS
	RED_WINS
	DRAW
)

type GameSettings struct {
	Rows    int
	Columns int
}

// Board size can be modified in the game overlay while not in game
type Power struct {
	Board     [][]rune
	IsPlaying Player
	Settings  GameSettings
	State     GameState
}

type Coordinate struct {
	Column int
	Row    int
}

func initBoard(settings GameSettings) [][]rune {
	board := make([][]rune, settings.Rows)

	for i := 0; i < settings.Rows; i++ {
		board[i] = make([]rune, settings.Columns)
	}
	return board
}

func NewGameInstance(settings GameSettings) *Power {
	return &Power{
		Board:     initBoard(settings),
		IsPlaying: BLUE,
		Settings:  settings,
		State:     ONGOING,
	}
}

func (p *Power) MakeMove(coord Coordinate) {
	// Don't allow moves if game is over
	if p.State != ONGOING {
		return
	}

	if coord.Column < 0 || coord.Column >= p.Settings.Columns {
		return
	}

	for row := p.Settings.Rows - 1; row >= 0; row-- {
		if p.Board[row][coord.Column] == 0 {
			if p.IsPlaying == BLUE {
				p.Board[row][coord.Column] = 'B'
			} else {
				p.Board[row][coord.Column] = 'R'
			}

			// Check for victory or draw after the move
			p.checkGameState(row, coord.Column)

			// Switch player only if game is still ongoing
			if p.State == ONGOING {
				if p.IsPlaying == BLUE {
					p.IsPlaying = RED
				} else {
					p.IsPlaying = BLUE
				}
			}
			break
		}
	}
}

// checkGameState checks for victory or draw conditions after a move
func (p *Power) checkGameState(lastRow, lastCol int) {
	// Check if current player won
	if p.checkVictory(lastRow, lastCol) {
		if p.IsPlaying == BLUE {
			p.State = BLUE_WINS
		} else {
			p.State = RED_WINS
		}
		return
	}

	// Check for draw (board full)
	if p.isBoardFull() {
		p.State = DRAW
	}
}

// checkVictory checks if the last move resulted in a victory
func (p *Power) checkVictory(row, col int) bool {
	piece := p.Board[row][col]

	// Check horizontal
	if p.checkDirection(row, col, 0, 1, piece) >= 4 {
		return true
	}

	// Check vertical
	if p.checkDirection(row, col, 1, 0, piece) >= 4 {
		return true
	}

	// Check diagonal (top-left to bottom-right)
	if p.checkDirection(row, col, 1, 1, piece) >= 4 {
		return true
	}

	// Check diagonal (top-right to bottom-left)
	if p.checkDirection(row, col, 1, -1, piece) >= 4 {
		return true
	}

	return false
}

// checkDirection counts consecutive pieces in a given direction
func (p *Power) checkDirection(row, col, deltaRow, deltaCol int, piece rune) int {
	count := 1 // Count the current piece

	// Check in positive direction
	r, c := row+deltaRow, col+deltaCol
	for r >= 0 && r < p.Settings.Rows && c >= 0 && c < p.Settings.Columns && p.Board[r][c] == piece {
		count++
		r += deltaRow
		c += deltaCol
	}

	// Check in negative direction
	r, c = row-deltaRow, col-deltaCol
	for r >= 0 && r < p.Settings.Rows && c >= 0 && c < p.Settings.Columns && p.Board[r][c] == piece {
		count++
		r -= deltaRow
		c -= deltaCol
	}

	return count
}

// isBoardFull checks if the board is completely full
func (p *Power) isBoardFull() bool {
	for col := 0; col < p.Settings.Columns; col++ {
		if p.Board[0][col] == 0 {
			return false
		}
	}
	return true
}

// GetGameState returns the current game state
func (p *Power) GetGameState() GameState {
	return p.State
}

// IsGameOver returns true if the game has ended
func (p *Power) IsGameOver() bool {
	return p.State != ONGOING
}

// GetWinner returns the winning player, or nil if no winner yet
func (p *Power) GetWinner() *Player {
	switch p.State {
	case BLUE_WINS:
		winner := BLUE
		return &winner
	case RED_WINS:
		winner := RED
		return &winner
	default:
		return nil
	}
}

// ResetGame resets the game to initial state
func (p *Power) ResetGame() {
	p.Board = initBoard(p.Settings)
	p.IsPlaying = BLUE
	p.State = ONGOING
}

// IsValidMove checks if a move is valid without making it
func (p *Power) IsValidMove(coord Coordinate) bool {
	if p.State != ONGOING {
		return false
	}

	if coord.Column < 0 || coord.Column >= p.Settings.Columns {
		return false
	}

	// Check if column has space (top row is empty)
	return p.Board[0][coord.Column] == 0
}

// String returns a string representation of the player
func (player Player) String() string {
	switch player {
	case BLUE:
		return "Blue"
	case RED:
		return "Red"
	default:
		return "Unknown"
	}
}

// String returns a string representation of the game state
func (state GameState) String() string {
	switch state {
	case ONGOING:
		return "Ongoing"
	case BLUE_WINS:
		return "Blue Wins"
	case RED_WINS:
		return "Red Wins"
	case DRAW:
		return "Draw"
	default:
		return "Unknown"
	}
}

// GetCurrentPlayer returns the player whose turn it is
func (p *Power) GetCurrentPlayer() Player {
	return p.IsPlaying
}
