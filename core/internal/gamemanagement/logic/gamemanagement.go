package logic

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/lithammer/shortuuid"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	gerr "github.com/yoktobit/secretsanta/internal/gamemanagement/logic/errors"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic/to"
	gda "github.com/yoktobit/secretsanta/internal/general/dataaccess"
	glogic "github.com/yoktobit/secretsanta/internal/general/logic"
	"golang.org/x/crypto/bcrypt"

	log "github.com/sirupsen/logrus"
)

// Gamemanagement contains the business logic for this app
type Gamemanagement interface {
	Connection() gda.Connection
	CreateNewGame(createGameTo to.CreateGameTo) (to.CreateGameResponseTo, error)
	AddPlayerToGame(addPlayerTo to.AddRemovePlayerTo) error
	RemovePlayerFromGame(removePlayerTo to.AddRemovePlayerTo) error
	RegisterPlayerPassword(registerPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) error
	AddException(addExceptionTo to.AddExceptionTo) error
	GetBasicGameByCode(code string) (to.GetBasicGameResponseTo, error)
	GetFullGameByCode(code string, playerName string) (to.GetFullGameResponseTo, error)
	GetPlayersByCode(code string) ([]to.PlayerResponseTo, error)
	GetPlayerRoleByCodeAndName(code string, name string) (string, error)
	GetExceptionsByCode(code string) ([]to.ExceptionResponseTo, error)
	LoginPlayer(loginPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) to.RegisterLoginPlayerPasswordResponseTo
	DrawGame(drawGameTo to.DrawGameTo) (to.DrawGameResponseTo, error)
	ResetGame(gameCode string) error
}

type gamemanagement struct {
	connection                gda.Connection
	gameRepository            dataaccess.GameRepository
	playerRepository          dataaccess.PlayerRepository
	playerExceptionRepository dataaccess.PlayerExceptionRepository
	random                    glogic.Randomizer
}

// NewGamemanagement is the factory method to create a new Gamemanagement
func NewGamemanagement(connection gda.Connection, gameRepository dataaccess.GameRepository, playerRepository dataaccess.PlayerRepository, playerExceptionRepository dataaccess.PlayerExceptionRepository, random glogic.Randomizer) Gamemanagement {

	return &gamemanagement{connection: connection, gameRepository: gameRepository, playerRepository: playerRepository, playerExceptionRepository: playerExceptionRepository, random: random}
}

// Connection returns the database connection
func (gamemanagement *gamemanagement) Connection() gda.Connection {

	return gamemanagement.connection
}

// CreateNewGame creates a new Game
func (gamemanagement *gamemanagement) CreateNewGame(createGameTo to.CreateGameTo) (to.CreateGameResponseTo, error) {
	err := validator.New().Struct(createGameTo)
	if err != nil {
		return to.CreateGameResponseTo{}, err
	}
	code := gamemanagement.generateCode()
	game := dataaccess.Game{Code: code, Title: createGameTo.Title, Description: createGameTo.Description, Status: dataaccess.StatusCreated.String()}
	hashedPassword := gamemanagement.generatePassword(createGameTo.AdminPassword)
	gamemanagement.Connection().NewTransaction(func(c gda.Connection) error {
		gamemanagement.gameRepository.CreateGame(c, &game)
		player := dataaccess.Player{Name: createGameTo.AdminUser, Password: hashedPassword, GameID: game.ID, Role: dataaccess.RoleAdmin.String(), Status: dataaccess.StatusReady.String()}
		gamemanagement.playerRepository.CreatePlayer(c, &player)
		return nil
	})
	result := to.CreateGameResponseTo{Code: code}
	return result, nil
}

// AddPlayerToGame adds a new player to an existing game
func (gamemanagement *gamemanagement) AddPlayerToGame(addPlayerTo to.AddRemovePlayerTo) error {
	err := validator.New().Struct(addPlayerTo)
	if err != nil {
		return err
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(addPlayerTo.GameCode)
	if err != nil {
		return err
	}
	player := dataaccess.Player{Name: addPlayerTo.Name, GameID: game.ID, Role: dataaccess.RolePlayer.String()}
	game.Status = dataaccess.StatusWaiting.String()
	gamemanagement.Connection().NewTransaction(func(c gda.Connection) error {
		gamemanagement.playerRepository.CreatePlayer(c, &player)
		gamemanagement.gameRepository.UpdateGame(c, &game)
		return nil
	})
	return nil
}

// RemovePlayerFromGame removes an existing player from an existing game
func (gamemanagement *gamemanagement) RemovePlayerFromGame(removePlayerTo to.AddRemovePlayerTo) error {
	err := validator.New().Struct(removePlayerTo)
	if err != nil {
		return err
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(removePlayerTo.GameCode)
	if err != nil {
		return err
	}
	player, err := gamemanagement.playerRepository.FindPlayerByNameAndGameID(removePlayerTo.Name, game.ID)
	if err == nil {
		gamemanagement.Connection().NewTransaction(func(c gda.Connection) error {
			gamemanagement.playerExceptionRepository.DeleteExceptionByPlayerID(c, player.ID)
			gamemanagement.playerRepository.DeletePlayerByNameAndGameID(c, removePlayerTo.Name, game.ID)
			gamemanagement.refreshGameStatus(c, &game)
			return nil
		})
	}
	return err
}

// RegisterPlayerPassword registers the password for a player and tells he/she is ready to go
func (gamemanagement *gamemanagement) RegisterPlayerPassword(registerPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) error {
	err := validator.New().Struct(registerPlayerPasswordTo)
	if err != nil {
		return err
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(registerPlayerPasswordTo.GameCode)
	if err != nil {
		return err
	}
	player, err := gamemanagement.playerRepository.FindPlayerByNameAndGameID(registerPlayerPasswordTo.Name, game.ID)
	if err != nil {
		return err
	}
	player.Password = gamemanagement.generatePassword(registerPlayerPasswordTo.Password)
	player.Status = dataaccess.StatusReady.String()
	gamemanagement.Connection().NewTransaction(func(c gda.Connection) error {
		gamemanagement.playerRepository.UpdatePlayer(c, &player)
		gamemanagement.refreshGameStatus(c, &game)
		return nil
	})
	return nil
}

// AddException adds a new exception so that PlayerA doesnt have to gift PlayerB
func (gamemanagement *gamemanagement) AddException(addExceptionTo to.AddExceptionTo) error {
	err := validator.New().Struct(addExceptionTo)
	if err != nil {
		return err
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(addExceptionTo.GameCode)
	if err != nil {
		return err
	}
	playerA, err := gamemanagement.playerRepository.FindPlayerByNameAndGameID(addExceptionTo.NameA, game.ID)
	if err != nil {
		return err
	}
	playerB, err := gamemanagement.playerRepository.FindPlayerByNameAndGameID(addExceptionTo.NameB, game.ID)
	if err != nil {
		return err
	}
	_, err = gamemanagement.playerExceptionRepository.FindExceptionByIds(playerA.ID, playerB.ID, game.ID)
	if err != nil {
		playerException := dataaccess.PlayerException{PlayerA: playerA, PlayerB: playerB, GameID: game.ID}
		gamemanagement.connection.NewTransaction(func(c gda.Connection) error {
			gamemanagement.playerExceptionRepository.CreatePlayerException(c, &playerException)
			return nil
		})
	} else {
		return gerr.ErrPlayerExceptionAlreadyExists
	}
	return nil
}

// GetBasicGameByCode fetches the game from the DB
func (gamemanagement *gamemanagement) GetBasicGameByCode(code string) (to.GetBasicGameResponseTo, error) {

	if code == "" {
		return to.GetBasicGameResponseTo{}, errors.New("Code must not be empty")
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(code)
	if err != nil {
		return to.GetBasicGameResponseTo{}, err
	}
	gameResponseTo := to.GetBasicGameResponseTo{Title: game.Title, Description: game.Description, Code: game.Code}
	return gameResponseTo, nil
}

// GetFullGameByCode fetches the game from the DB
func (gamemanagement *gamemanagement) GetFullGameByCode(code string, playerName string) (to.GetFullGameResponseTo, error) {
	if code == "" {
		return to.GetFullGameResponseTo{}, errors.New("Code must not be empty")
	}
	if playerName == "" {
		return to.GetFullGameResponseTo{}, errors.New("playerName must not be empty")
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(code)
	if err != nil {
		return to.GetFullGameResponseTo{}, err
	}
	gameResponseTo := to.GetFullGameResponseTo{Title: game.Title, Description: game.Description, Status: game.Status, Code: game.Code}
	if game.Status == dataaccess.StatusDrawn.String() {
		player, err := gamemanagement.playerRepository.FindPlayerWithAssociationsByNameAndGameID(playerName, game.ID)
		if err != nil {
			return to.GetFullGameResponseTo{}, errors.New("Player not found")
		}
		gameResponseTo.Gifted = player.Gifted.Name
	}
	return gameResponseTo, nil
}

// GetPlayersByCode fetches the players of a game from the DB
func (gamemanagement *gamemanagement) GetPlayersByCode(code string) ([]to.PlayerResponseTo, error) {
	if code == "" {
		return make([]to.PlayerResponseTo, 0), errors.New("Code must not be empty")
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(code)
	if err != nil {
		return make([]to.PlayerResponseTo, 0), err
	}
	players, err := gamemanagement.playerRepository.FindPlayersByGameID(game.ID)
	if err != nil {
		return make([]to.PlayerResponseTo, 0), err
	}
	playerResponseTos := make([]to.PlayerResponseTo, 0)
	for _, player := range players {
		playerResponseTo := to.PlayerResponseTo{Name: player.Name, Status: player.Status}
		playerResponseTos = append(playerResponseTos, playerResponseTo)
	}
	return playerResponseTos, nil
}

// GetPlayerRoleByCodeAndName fetches the player role in a game from the DB
func (gamemanagement *gamemanagement) GetPlayerRoleByCodeAndName(code string, name string) (string, error) {
	if code == "" {
		return "", errors.New("Code must not be empty")
	}
	if name == "" {
		return "", errors.New("Name must not be empty")
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(code)
	if err != nil {
		return "", err
	}
	player, err := gamemanagement.playerRepository.FindPlayerByNameAndGameID(name, game.ID)
	if err == nil {
		log.WithField("role", player.Role).Debug("Rolle gefunden")
		return player.Role, nil
	}
	return "", err
}

// GetExceptionsByCode returns the draw exceptions in a game
func (gamemanagement *gamemanagement) GetExceptionsByCode(code string) ([]to.ExceptionResponseTo, error) {
	exceptionResponseTos := make([]to.ExceptionResponseTo, 0)
	if code == "" {
		return exceptionResponseTos, errors.New("Code must not be empty")
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(code)
	if err != nil {
		return exceptionResponseTos, err
	}
	playerExceptions, err := gamemanagement.playerExceptionRepository.FindExceptionsWithAssociationsByGameID(game.ID)
	if err != nil {
		return exceptionResponseTos, err
	}
	for _, playerException := range playerExceptions {
		exceptionResponseTo := to.ExceptionResponseTo{NameA: playerException.PlayerA.Name, NameB: playerException.PlayerB.Name}
		exceptionResponseTos = append(exceptionResponseTos, exceptionResponseTo)
	}
	return exceptionResponseTos, err
}

// LoginPlayer logs in a player
func (gamemanagement *gamemanagement) LoginPlayer(loginPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) to.RegisterLoginPlayerPasswordResponseTo {
	var player dataaccess.Player
	var loginPlayerPasswordResponseTo to.RegisterLoginPlayerPasswordResponseTo
	err := validator.New().Struct(loginPlayerPasswordTo)
	if err != nil {
		gamemanagement.writeLoginError(&loginPlayerPasswordResponseTo)
		return loginPlayerPasswordResponseTo
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(loginPlayerPasswordTo.GameCode)
	if err != nil {
		gamemanagement.writeLoginError(&loginPlayerPasswordResponseTo)
		return loginPlayerPasswordResponseTo
	}
	player, err = gamemanagement.playerRepository.FindPlayerByNameAndGameID(loginPlayerPasswordTo.Name, game.ID)
	if err != nil {
		gamemanagement.writeLoginError(&loginPlayerPasswordResponseTo)
		return loginPlayerPasswordResponseTo
	}
	if player.Password == "" {
		gamemanagement.RegisterPlayerPassword(loginPlayerPasswordTo)
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(loginPlayerPasswordTo.Password))
		if err != nil {
			gamemanagement.writeLoginError(&loginPlayerPasswordResponseTo)
			return loginPlayerPasswordResponseTo
		}
	}
	loginPlayerPasswordResponseTo.Ok = true
	return loginPlayerPasswordResponseTo
}

// DrawGame draws the lots
func (gamemanagement *gamemanagement) DrawGame(drawGameTo to.DrawGameTo) (to.DrawGameResponseTo, error) {
	if drawGameTo.GameCode == "" {
		return to.DrawGameResponseTo{}, errors.New("GameCode must not be empty")
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(drawGameTo.GameCode)
	if err != nil {
		return to.DrawGameResponseTo{}, err
	}
	players, err := gamemanagement.playerRepository.FindPlayersByGameID(game.ID)
	if err != nil {
		return to.DrawGameResponseTo{}, err
	}
	exceptions, err := gamemanagement.playerExceptionRepository.FindExceptionsWithAssociationsByGameID(game.ID)
	if err != nil {
		return to.DrawGameResponseTo{}, err
	}
	var lots map[*dataaccess.Player]*dataaccess.Player = make(map[*dataaccess.Player]*dataaccess.Player)
	tries := 0
	ok := false
	for !ok && tries < 100 {
		remaining := make([]*dataaccess.Player, len(players))
		copy(remaining, players)
		for _, player := range players {
			randomNumber := gamemanagement.random.NextInt(len(remaining))
			lots[player] = remaining[randomNumber]
			log.WithFields(log.Fields{"giftee": player.Name, "gifted": lots[player].Name}).Debug("Los")
			remaining = append(remaining[:randomNumber], remaining[randomNumber+1:]...)
			log.WithField("remaining", len(remaining)).Debug("Remaining")
			singleOk := gamemanagement.isAllowed(exceptions, player, lots[player])
			if !singleOk {
				break
			}
		}
		ok = gamemanagement.checkResult(exceptions, lots)
		tries++
	}
	drawGameResponseTo := to.DrawGameResponseTo{}
	if ok {
		game.Status = dataaccess.StatusDrawn.String()
		gamemanagement.Connection().NewTransaction(func(c gda.Connection) error {
			gamemanagement.saveLots(c, lots)
			gamemanagement.gameRepository.UpdateGame(c, &game)
			return nil
		})
	} else {
		log.Warn("Keine plausible Auslosung gefunden")
		drawGameResponseTo.Message = "Nach 100 Versuchen wurde kein plausibles Ergebnis gefunden. Bitte nochmal versuchen oder weniger Ausnahmen definieren."
	}
	drawGameResponseTo.Ok = ok
	return drawGameResponseTo, nil
}

// ResetGame resets a game
func (gamemanagement *gamemanagement) ResetGame(code string) error {
	if code == "" {
		return errors.New("Code must not be empty")
	}
	game, err := gamemanagement.gameRepository.FindGameByCode(code)
	if err != nil {
		return err
	}
	game.Status = dataaccess.StatusReady.String()
	gamemanagement.Connection().NewTransaction(func(c gda.Connection) error {
		gamemanagement.gameRepository.UpdateGame(c, &game)
		return nil
	})
	return nil
}

func (gamemanagement *gamemanagement) saveLots(c gda.Connection, lots map[*dataaccess.Player]*dataaccess.Player) {

	for giftee, gifted := range lots {

		giftee.GiftedID = &gifted.ID
		gamemanagement.playerRepository.UpdatePlayer(c, giftee)
	}
}

func (gamemanagement *gamemanagement) checkResult(exceptions []*dataaccess.PlayerException, lots map[*dataaccess.Player]*dataaccess.Player) bool {

	for giftee, gifted := range lots {
		if !gamemanagement.isAllowed(exceptions, giftee, gifted) {
			return false
		}
	}
	return true
}

func (gamemanagement *gamemanagement) isAllowed(exceptions []*dataaccess.PlayerException, giftee *dataaccess.Player, gifted *dataaccess.Player) bool {

	if giftee.ID == gifted.ID {
		return false
	}
	for _, playerException := range exceptions {
		if playerException.PlayerA.ID == giftee.ID && playerException.PlayerB.ID == gifted.ID {
			return false
		}
	}
	return true
}

func (gamemanagement *gamemanagement) writeLoginError(loginPlayerPasswordResponseTo *to.RegisterLoginPlayerPasswordResponseTo) {
	loginPlayerPasswordResponseTo.Message = "Falsche Game-ID, falscher Nutzername oder falsches Passwort"
	loginPlayerPasswordResponseTo.Ok = false
}

func (gamemanagement *gamemanagement) generateCode() string {
	return shortuuid.New()
}

func (gamemanagement *gamemanagement) generatePassword(plainPassword string) string {
	plainPasswordByte := []byte(plainPassword)
	// I see no way how this may fail in real life
	hash, _ := bcrypt.GenerateFromPassword(plainPasswordByte, bcrypt.DefaultCost)
	return string(hash)
}

func (gamemanagement *gamemanagement) refreshGameStatus(c gda.Connection, game *dataaccess.Game) error {
	_, exists, err := gamemanagement.playerRepository.FindFirstUnreadyPlayerByGameID(game.ID)
	if err != nil {
		return err
	}
	if !exists {
		game.Status = dataaccess.StatusReady.String()
		gamemanagement.gameRepository.UpdateGame(c, game)
	}
	return err
}
