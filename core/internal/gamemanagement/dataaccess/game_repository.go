package dataaccess

import (
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MigrateDb migrates the DB to the current schema version
func MigrateDb(database *gorm.DB) {

	database.AutoMigrate(&Game{})
	database.AutoMigrate(&Player{})
	database.AutoMigrate(&PlayerException{})
}

// CreateGame creates a game
func CreateGame(game *Game) {

	dataaccess.DB.Create(game)
}

// CreatePlayer creates a player
func CreatePlayer(player *Player) {

	dataaccess.DB.Create(player)
}

// CreatePlayerException creates an Exception
func CreatePlayerException(playerException *PlayerException) {

	dataaccess.DB.Create(playerException)
}

// GetGameByCode receives a game by code
func GetGameByCode(code string) (Game, error) {

	var game Game
	result := dataaccess.DB.First(&game, "code = ?", code)
	return game, result.Error
}

// UpdateGame updates a game
func UpdateGame(game *Game) {

	dataaccess.DB.Save(game)
}

// UpdatePlayer updates a player
func UpdatePlayer(player *Player) {

	dataaccess.DB.Save(player)
}

// GetPlayerByNameAndGameID Get Player by name and game id
func GetPlayerByNameAndGameID(name string, gameID uint) (Player, error) {

	var player Player
	result := dataaccess.DB.First(&player, "name = ? AND game_id = ?", name, gameID)
	return player, result.Error
}

// GetExceptionByIds receives an exception by player ids and game id
func GetExceptionByIds(playerAId uint, playerBId uint, gameID uint) (PlayerException, error) {

	var existingException PlayerException
	result := dataaccess.DB.First(&existingException, "player_a_id = ? AND player_b_id = ? AND game_id = ?", playerAId, playerBId, gameID)
	return existingException, result.Error
}

// GetPlayerWithAssociationsByNameAndGameID Get a Player By Name and Game ID including Associations
func GetPlayerWithAssociationsByNameAndGameID(playerName string, gameID uint) Player {

	var player Player
	dataaccess.DB.Preload(clause.Associations).First(&player, "name = ? AND game_id = ?", playerName, gameID)
	return player
}

// GetPlayersByGameID Get all Players by Game ID
func GetPlayersByGameID(gameID uint) []*Player {

	var players []*Player
	dataaccess.DB.Where("game_id = ?", gameID).Find(&players)
	return players
}

// GetExceptionsWithAssociationsByGameID Get existing Exceptions by Game ID including Associations
func GetExceptionsWithAssociationsByGameID(gameID uint) []*PlayerException {

	var playerExceptions []*PlayerException
	dataaccess.DB.Where("game_id = ?", gameID).Preload(clause.Associations).Find(&playerExceptions)
	return playerExceptions
}

// GetFirstUnreadyPlayerByGameID Get the first unready Player for a GameId
func GetFirstUnreadyPlayerByGameID(gameID uint) (Player, error) {

	var otherPlayer Player
	result := dataaccess.DB.First(&otherPlayer, "game_id = ? AND status != ?", gameID, Ready.String())
	return otherPlayer, result.Error
}

// DeleteExceptionByPlayerID deletes the game by the IDs of Player A and Player B
func DeleteExceptionByPlayerID(playerID uint) {

	var exception PlayerException
	dataaccess.DB.Delete(&exception, "player_a_id = ? OR player_b_id = ?", playerID, playerID)
}

// DeletePlayerByNameAndGameID deletes a player by name and game ID
func DeletePlayerByNameAndGameID(playerName string, gameID uint) {

	var player Player
	dataaccess.DB.Delete(&player, "name = ? AND game_id = ?", playerName, gameID)
}
