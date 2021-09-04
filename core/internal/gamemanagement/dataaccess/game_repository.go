package dataaccess

import (
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
	"gorm.io/gorm"
)

// GameRepository holds all the database access functions
type GameRepository interface {
	CreateGame(c dataaccess.Connection, game *Game)
	UpdateGame(c dataaccess.Connection, game *Game)
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
func (gameRepository *gameRepository) CreateGame(c dataaccess.Connection, game *Game) {

	c.Connection().Create(game)
}

// FindGameByCode receives a game by code
func (gameRepository *gameRepository) FindGameByCode(code string) (Game, error) {

	var game Game
	result := gameRepository.connection.Connection().Where("code = ?", code).Limit(1).Find(&game)
	if result.RowsAffected == 0 {
		return game, gorm.ErrRecordNotFound
	}
	return game, result.Error
}

// UpdateGame updates a game
func (gameRepository *gameRepository) UpdateGame(c dataaccess.Connection, game *Game) {

	c.Connection().Save(game)
}
