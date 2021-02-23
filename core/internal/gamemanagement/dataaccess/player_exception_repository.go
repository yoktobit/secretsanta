package dataaccess

import (
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PlayerExceptionRepository holds all the database access functions
type PlayerExceptionRepository interface {
	CreatePlayerException(c dataaccess.Connection, playerException *PlayerException)
	DeleteExceptionByPlayerID(c dataaccess.Connection, playerID uint)
	FindExceptionByIds(playerAId uint, playerBId uint, gameID uint) (PlayerException, error)
	FindExceptionsWithAssociationsByGameID(gameID uint) ([]*PlayerException, error)
}

type playerExceptionRepository struct {
	connection dataaccess.Connection
}

// NewPlayerExceptionRepository is the factory method for creating a PlayerException repository
func NewPlayerExceptionRepository(connection dataaccess.Connection) PlayerExceptionRepository {

	return &playerExceptionRepository{connection: connection}
}

// CreatePlayerException creates an Exception
func (playerExceptionRepository *playerExceptionRepository) CreatePlayerException(c dataaccess.Connection, playerException *PlayerException) {

	c.Connection().Create(playerException)
}

// FindExceptionByIds receives an exception by player ids and game id
func (playerExceptionRepository *playerExceptionRepository) FindExceptionByIds(playerAId uint, playerBId uint, gameID uint) (PlayerException, error) {

	var existingException PlayerException
	result := playerExceptionRepository.connection.Connection().Where("player_a_id = ? AND player_b_id = ? AND game_id = ?", playerAId, playerBId, gameID).Limit(1).Find(&existingException)
	if result.RowsAffected == 0 {
		return existingException, gorm.ErrRecordNotFound
	}
	return existingException, nil
}

// FindExceptionsWithAssociationsByGameID Get existing Exceptions by Game ID including Associations
func (playerExceptionRepository *playerExceptionRepository) FindExceptionsWithAssociationsByGameID(gameID uint) ([]*PlayerException, error) {

	var playerExceptions []*PlayerException
	result := playerExceptionRepository.connection.Connection().Where("game_id = ?", gameID).Preload(clause.Associations).Find(&playerExceptions)
	return playerExceptions, result.Error
}

// DeleteExceptionByPlayerID deletes the game by the IDs of Player A and Player B
func (playerExceptionRepository *playerExceptionRepository) DeleteExceptionByPlayerID(c dataaccess.Connection, playerID uint) {

	var exception PlayerException
	c.Connection().Delete(&exception, "player_a_id = ? OR player_b_id = ?", playerID, playerID)
}
