package dataaccess

import (
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PlayerRepository holds all the database access functions
type PlayerRepository interface {
	CreatePlayer(c dataaccess.Connection, player *Player)
	UpdatePlayer(c dataaccess.Connection, player *Player)
	DeletePlayerByNameAndGameID(c dataaccess.Connection, playerName string, gameID uint)
	FindPlayerByNameAndGameID(name string, gameID uint) (Player, error)
	FindPlayerWithAssociationsByNameAndGameID(playerName string, gameID uint) Player
	FindFirstUnreadyPlayerByGameID(gameID uint) (Player, error)
	FindPlayersByGameID(gameID uint) []*Player
}

type playerRepository struct {
	connection dataaccess.Connection
}

// NewPlayerRepository is the factory method for creating a player repository
func NewPlayerRepository(connection dataaccess.Connection) PlayerRepository {

	return &playerRepository{connection: connection}
}

// CreatePlayer creates a player
func (playerRepository *playerRepository) CreatePlayer(c dataaccess.Connection, player *Player) {

	c.Connection().Create(player)
}

// UpdatePlayer updates a player
func (playerRepository *playerRepository) UpdatePlayer(c dataaccess.Connection, player *Player) {

	c.Connection().Save(player)
}

// FindPlayerByNameAndGameID Get Player by name and game id
func (playerRepository *playerRepository) FindPlayerByNameAndGameID(name string, gameID uint) (Player, error) {

	var player Player
	result := playerRepository.connection.Connection().Where("name = ? AND game_id = ?", name, gameID).Limit(1).Find(&player)
	if result.RowsAffected == 0 {
		return player, gorm.ErrRecordNotFound
	}
	return player, result.Error
}

// FindPlayerWithAssociationsByNameAndGameID Get a Player By Name and Game ID including Associations
func (playerRepository *playerRepository) FindPlayerWithAssociationsByNameAndGameID(playerName string, gameID uint) Player {

	var player Player
	playerRepository.connection.Connection().Preload(clause.Associations).First(&player, "name = ? AND game_id = ?", playerName, gameID)
	return player
}

// FindPlayersByGameID Get all Players by Game ID
func (playerRepository *playerRepository) FindPlayersByGameID(gameID uint) []*Player {

	var players []*Player
	playerRepository.connection.Connection().Where("game_id = ?", gameID).Find(&players)
	return players
}

// FindFirstUnreadyPlayerByGameID Get the first unready Player for a GameId
func (playerRepository *playerRepository) FindFirstUnreadyPlayerByGameID(gameID uint) (Player, error) {

	var otherPlayer Player
	result := playerRepository.connection.Connection().First(&otherPlayer, "game_id = ? AND status != ?", gameID, StatusReady.String())
	return otherPlayer, result.Error
}

// DeletePlayerByNameAndGameID deletes a player by name and game ID
func (playerRepository *playerRepository) DeletePlayerByNameAndGameID(c dataaccess.Connection, playerName string, gameID uint) {

	var player Player
	c.Connection().Delete(&player, "name = ? AND game_id = ?", playerName, gameID)
}
