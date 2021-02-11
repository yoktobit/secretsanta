package dataaccess

import (
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

// GameRepository holds all the database access functions
type GameRepository interface {
	CreateGame(game *Game)
	UpdateGame(game *Game)
	FindGameByCode(code string) (Game, error)
}

type gameRepository struct {
	connection dataaccess.Connection
}

// NewGameRepository is the factory method for creating a game repository
func NewGameRepository(connection dataaccess.Connection) GameRepository {

	return &gameRepository{connection: connection}
}

// CreateGame creates a game
func (gameRepository *gameRepository) CreateGame(game *Game) {

	gameRepository.connection.Connection().Create(game)
	gameRepository.connection.Connection().Commit()
}

// FindGameByCode receives a game by code
func (gameRepository *gameRepository) FindGameByCode(code string) (Game, error) {

	var game Game
	result := gameRepository.connection.Connection().First(&game, "code = ?", code)
	return game, result.Error
}

// UpdateGame updates a game
func (gameRepository *gameRepository) UpdateGame(game *Game) {

	gameRepository.connection.Connection().Save(game)
}
