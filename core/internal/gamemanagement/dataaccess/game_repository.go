package dataaccess

import (
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GameRepository holds all the database access functions
type GameRepository interface {
	MigrateDb(database *gorm.DB)
	CreateGame(game *Game)
	CreatePlayer(player *Player)
	CreatePlayerException(playerException *PlayerException)
	GetGameByCode(code string) (Game, error)
	UpdateGame(game *Game)
	UpdatePlayer(player *Player)
	GetPlayerByNameAndGameID(name string, gameID uint) (Player, error)
	GetExceptionByIds(playerAId uint, playerBId uint, gameID uint) (PlayerException, error)
	GetPlayerWithAssociationsByNameAndGameID(playerName string, gameID uint) Player
	GetPlayersByGameID(gameID uint) []*Player
	GetExceptionsWithAssociationsByGameID(gameID uint) []*PlayerException
	GetFirstUnreadyPlayerByGameID(gameID uint) (Player, error)
	DeleteExceptionByPlayerID(playerID uint)
	DeletePlayerByNameAndGameID(playerName string, gameID uint)
}

type gameRepository struct {
	connection dataaccess.Connection
}

// NewGameRepository is the factory method for creating a game repository
func NewGameRepository(connection dataaccess.Connection) GameRepository {

	return &gameRepository{connection: connection}
}

// MigrateDb migrates the DB to the current schema version
func (gameRepository *gameRepository) MigrateDb(database *gorm.DB) {

	gameRepository.connection.Connection().AutoMigrate(&Game{})
	gameRepository.connection.Connection().AutoMigrate(&Player{})
	gameRepository.connection.Connection().AutoMigrate(&PlayerException{})
}

// CreateGame creates a game
func (gameRepository *gameRepository) CreateGame(game *Game) {

	gameRepository.connection.Connection().Create(game)
}

// CreatePlayer creates a player
func (gameRepository *gameRepository) CreatePlayer(player *Player) {

	gameRepository.connection.Connection().Create(player)
}

// CreatePlayerException creates an Exception
func (gameRepository *gameRepository) CreatePlayerException(playerException *PlayerException) {

	gameRepository.connection.Connection().Create(playerException)
}

// GetGameByCode receives a game by code
func (gameRepository *gameRepository) GetGameByCode(code string) (Game, error) {

	var game Game
	result := gameRepository.connection.Connection().First(&game, "code = ?", code)
	return game, result.Error
}

// UpdateGame updates a game
func (gameRepository *gameRepository) UpdateGame(game *Game) {

	gameRepository.connection.Connection().Save(game)
}

// UpdatePlayer updates a player
func (gameRepository *gameRepository) UpdatePlayer(player *Player) {

	gameRepository.connection.Connection().Save(player)
}

// GetPlayerByNameAndGameID Get Player by name and game id
func (gameRepository *gameRepository) GetPlayerByNameAndGameID(name string, gameID uint) (Player, error) {

	var player Player
	result := gameRepository.connection.Connection().First(&player, "name = ? AND game_id = ?", name, gameID)
	return player, result.Error
}

// GetExceptionByIds receives an exception by player ids and game id
func (gameRepository *gameRepository) GetExceptionByIds(playerAId uint, playerBId uint, gameID uint) (PlayerException, error) {

	var existingException PlayerException
	result := gameRepository.connection.Connection().First(&existingException, "player_a_id = ? AND player_b_id = ? AND game_id = ?", playerAId, playerBId, gameID)
	return existingException, result.Error
}

// GetPlayerWithAssociationsByNameAndGameID Get a Player By Name and Game ID including Associations
func (gameRepository *gameRepository) GetPlayerWithAssociationsByNameAndGameID(playerName string, gameID uint) Player {

	var player Player
	gameRepository.connection.Connection().Preload(clause.Associations).First(&player, "name = ? AND game_id = ?", playerName, gameID)
	return player
}

// GetPlayersByGameID Get all Players by Game ID
func (gameRepository *gameRepository) GetPlayersByGameID(gameID uint) []*Player {

	var players []*Player
	gameRepository.connection.Connection().Where("game_id = ?", gameID).Find(&players)
	return players
}

// GetExceptionsWithAssociationsByGameID Get existing Exceptions by Game ID including Associations
func (gameRepository *gameRepository) GetExceptionsWithAssociationsByGameID(gameID uint) []*PlayerException {

	var playerExceptions []*PlayerException
	gameRepository.connection.Connection().Where("game_id = ?", gameID).Preload(clause.Associations).Find(&playerExceptions)
	return playerExceptions
}

// GetFirstUnreadyPlayerByGameID Get the first unready Player for a GameId
func (gameRepository *gameRepository) GetFirstUnreadyPlayerByGameID(gameID uint) (Player, error) {

	var otherPlayer Player
	result := gameRepository.connection.Connection().First(&otherPlayer, "game_id = ? AND status != ?", gameID, Ready.String())
	return otherPlayer, result.Error
}

// DeleteExceptionByPlayerID deletes the game by the IDs of Player A and Player B
func (gameRepository *gameRepository) DeleteExceptionByPlayerID(playerID uint) {

	var exception PlayerException
	gameRepository.connection.Connection().Delete(&exception, "player_a_id = ? OR player_b_id = ?", playerID, playerID)
}

// DeletePlayerByNameAndGameID deletes a player by name and game ID
func (gameRepository *gameRepository) DeletePlayerByNameAndGameID(playerName string, gameID uint) {

	var player Player
	gameRepository.connection.Connection().Delete(&player, "name = ? AND game_id = ?", playerName, gameID)
}
