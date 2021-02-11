package logic

import (
	"math/rand"
	"time"

	"github.com/lithammer/shortuuid"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic/to"
	"golang.org/x/crypto/bcrypt"

	log "github.com/sirupsen/logrus"
)

// Gamemanagement contains the business logic for this app
type Gamemanagement interface {
	CreateNewGame(createGameTo to.CreateGameTo) to.CreateGameResponseTo
	AddPlayerToGame(addPlayerTo to.AddRemovePlayerTo)
	RemovePlayerFromGame(removePlayerTo to.AddRemovePlayerTo)
	RegisterPlayerPassword(registerPlayerPasswordTo to.RegisterLoginPlayerPasswordTo)
	AddException(addExceptionTo to.AddExceptionTo)
	GetBasicGameByCode(code string) to.GetBasicGameResponseTo
	GetFullGameByCode(code string, playerName string) to.GetFullGameResponseTo
	GetPlayersByCode(code string) []to.PlayerResponseTo
	GetPlayerRoleByCodeAndName(code string, name string) string
	GetExceptionsByCode(code string) []to.ExceptionResponseTo
	LoginPlayer(loginPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) to.RegisterLoginPlayerPasswordResponseTo
	DrawGame(drawGameTo to.DrawGameTo) to.DrawGameResponseTo
	ResetGame(gameCode string)
}

type gamemanagement struct {
	gameRepository            dataaccess.GameRepository
	playerRepository          dataaccess.PlayerRepository
	playerExceptionRepository dataaccess.PlayerExceptionRepository
}

// NewGamemanagement is the factory method to create a new Gamemanagement
func NewGamemanagement(gameRepository dataaccess.GameRepository, playerRepository dataaccess.PlayerRepository, playerExceptionRepository dataaccess.PlayerExceptionRepository) Gamemanagement {

	return &gamemanagement{gameRepository: gameRepository, playerRepository: playerRepository, playerExceptionRepository: playerExceptionRepository}
}

// CreateNewGame creates a new Game
func (gamemanagement *gamemanagement) CreateNewGame(createGameTo to.CreateGameTo) to.CreateGameResponseTo {
	code := gamemanagement.generateCode()
	game := dataaccess.Game{Code: code, Title: createGameTo.Title, Description: createGameTo.Description, Status: dataaccess.Created.String()}
	gamemanagement.gameRepository.CreateGame(&game)
	hashedPassword := gamemanagement.generatePassword(createGameTo.AdminPassword)
	player := dataaccess.Player{Name: createGameTo.AdminUser, Password: hashedPassword, GameID: game.ID, Role: dataaccess.RoleAdmin.String(), Status: dataaccess.Ready.String()}
	gamemanagement.playerRepository.CreatePlayer(&player)
	result := to.CreateGameResponseTo{Code: code, Link: gamemanagement.generateLink(code)}
	return result
}

// AddPlayerToGame adds a new player to an existing game
func (gamemanagement *gamemanagement) AddPlayerToGame(addPlayerTo to.AddRemovePlayerTo) {
	game, _ := gamemanagement.gameRepository.FindGameByCode(addPlayerTo.GameCode)
	player := dataaccess.Player{Name: addPlayerTo.Name, GameID: game.ID, Role: dataaccess.RolePlayer.String()}
	gamemanagement.playerRepository.CreatePlayer(&player)
	game.Status = dataaccess.Waiting.String()
	gamemanagement.gameRepository.UpdateGame(&game)
}

// RemovePlayerFromGame removes an existing player from an existing game
func (gamemanagement *gamemanagement) RemovePlayerFromGame(removePlayerTo to.AddRemovePlayerTo) {
	game, _ := gamemanagement.gameRepository.FindGameByCode(removePlayerTo.GameCode)
	player, error := gamemanagement.playerRepository.FindPlayerByNameAndGameID(removePlayerTo.Name, game.ID)
	if error == nil {
		gamemanagement.playerExceptionRepository.DeleteExceptionByPlayerID(player.ID)
		gamemanagement.playerRepository.DeletePlayerByNameAndGameID(removePlayerTo.Name, game.ID)
	}
	gamemanagement.refreshGameStatus(&game)
}

// RegisterPlayerPassword registers the password for a player and tells he/she is ready to go
func (gamemanagement *gamemanagement) RegisterPlayerPassword(registerPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) {
	game, _ := gamemanagement.gameRepository.FindGameByCode(registerPlayerPasswordTo.GameCode)
	player, _ := gamemanagement.playerRepository.FindPlayerByNameAndGameID(registerPlayerPasswordTo.Name, game.ID)
	player.Password = gamemanagement.generatePassword(registerPlayerPasswordTo.Password)
	player.Status = dataaccess.Ready.String()
	gamemanagement.playerRepository.UpdatePlayer(&player)
	gamemanagement.refreshGameStatus(&game)
}

// AddException adds a new exception so that PlayerA doesnt have to gift PlayerB
func (gamemanagement *gamemanagement) AddException(addExceptionTo to.AddExceptionTo) {
	game, _ := gamemanagement.gameRepository.FindGameByCode(addExceptionTo.GameCode)
	playerA, _ := gamemanagement.playerRepository.FindPlayerByNameAndGameID(addExceptionTo.NameA, game.ID)
	playerB, _ := gamemanagement.playerRepository.FindPlayerByNameAndGameID(addExceptionTo.NameB, game.ID)
	_, error := gamemanagement.playerExceptionRepository.FindExceptionByIds(playerA.ID, playerB.ID, game.ID)
	if error != nil {
		playerException := dataaccess.PlayerException{PlayerA: playerA, PlayerB: playerB, GameID: game.ID}
		gamemanagement.playerExceptionRepository.CreatePlayerException(&playerException)
	}
}

// GetBasicGameByCode fetches the game from the DB
func (gamemanagement *gamemanagement) GetBasicGameByCode(code string) to.GetBasicGameResponseTo {
	game, _ := gamemanagement.gameRepository.FindGameByCode(code)
	gameResponseTo := to.GetBasicGameResponseTo{Title: game.Title, Description: game.Description, Code: game.Code}
	return gameResponseTo
}

// GetFullGameByCode fetches the game from the DB
func (gamemanagement *gamemanagement) GetFullGameByCode(code string, playerName string) to.GetFullGameResponseTo {
	game, _ := gamemanagement.gameRepository.FindGameByCode(code)
	gameResponseTo := to.GetFullGameResponseTo{Title: game.Title, Description: game.Description, Status: game.Status, Code: game.Code}
	if game.Status == dataaccess.Drawn.String() {
		player := gamemanagement.playerRepository.FindPlayerWithAssociationsByNameAndGameID(playerName, game.ID)
		gameResponseTo.Gifted = player.Gifted.Name
	}
	return gameResponseTo
}

// GetPlayersByCode fetches the players of a game from the DB
func (gamemanagement *gamemanagement) GetPlayersByCode(code string) []to.PlayerResponseTo {
	game, _ := gamemanagement.gameRepository.FindGameByCode(code)
	players := gamemanagement.playerRepository.FindPlayersByGameID(game.ID)
	playerResponseTos := make([]to.PlayerResponseTo, 0)
	for _, player := range players {
		playerResponseTo := to.PlayerResponseTo{Name: player.Name, Status: player.Status}
		playerResponseTos = append(playerResponseTos, playerResponseTo)
	}
	return playerResponseTos
}

// GetPlayerRoleByCodeAndName fetches the player role in a game from the DB
func (gamemanagement *gamemanagement) GetPlayerRoleByCodeAndName(code string, name string) string {
	game, _ := gamemanagement.gameRepository.FindGameByCode(code)
	player, error := gamemanagement.playerRepository.FindPlayerByNameAndGameID(name, game.ID)
	if error == nil {
		log.WithField("role", player.Role).Debug("Rolle gefunden")
		return player.Role
	}
	return ""
}

// GetExceptionsByCode returns the draw exceptions in a game
func (gamemanagement *gamemanagement) GetExceptionsByCode(code string) []to.ExceptionResponseTo {
	game, _ := gamemanagement.gameRepository.FindGameByCode(code)
	playerExceptions := gamemanagement.playerExceptionRepository.FindExceptionsWithAssociationsByGameID(game.ID)
	exceptionResponseTos := make([]to.ExceptionResponseTo, 0)
	for _, playerException := range playerExceptions {
		exceptionResponseTo := to.ExceptionResponseTo{NameA: playerException.PlayerA.Name, NameB: playerException.PlayerB.Name}
		exceptionResponseTos = append(exceptionResponseTos, exceptionResponseTo)
	}
	return exceptionResponseTos
}

// LoginPlayer logs in a player
func (gamemanagement *gamemanagement) LoginPlayer(loginPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) to.RegisterLoginPlayerPasswordResponseTo {
	var player dataaccess.Player
	var loginPlayerPasswordResponseTo to.RegisterLoginPlayerPasswordResponseTo
	game, error := gamemanagement.gameRepository.FindGameByCode(loginPlayerPasswordTo.GameCode)
	if error != nil {
		log.WithField("gameCode", loginPlayerPasswordTo.GameCode).Info("Game not found")
		gamemanagement.writeLoginError(&loginPlayerPasswordResponseTo)
		return loginPlayerPasswordResponseTo
	}
	player, err := gamemanagement.playerRepository.FindPlayerByNameAndGameID(loginPlayerPasswordTo.Name, game.ID)
	if err != nil {
		log.WithFields(log.Fields{"name": loginPlayerPasswordTo.Name, "gameCode": loginPlayerPasswordTo.GameCode}).Info("Player not found")
		gamemanagement.writeLoginError(&loginPlayerPasswordResponseTo)
		return loginPlayerPasswordResponseTo
	}
	if player.Password == "" {
		gamemanagement.RegisterPlayerPassword(loginPlayerPasswordTo)
	} else {
		error := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(loginPlayerPasswordTo.Password))
		if error != nil {
			gamemanagement.writeLoginError(&loginPlayerPasswordResponseTo)
			return loginPlayerPasswordResponseTo
		}
	}
	loginPlayerPasswordResponseTo.Ok = true
	return loginPlayerPasswordResponseTo
}

// DrawGame draws the lots
func (gamemanagement *gamemanagement) DrawGame(drawGameTo to.DrawGameTo) to.DrawGameResponseTo {
	game, _ := gamemanagement.gameRepository.FindGameByCode(drawGameTo.GameCode)
	players := gamemanagement.playerRepository.FindPlayersByGameID(game.ID)
	exceptions := gamemanagement.playerExceptionRepository.FindExceptionsWithAssociationsByGameID(game.ID)
	var lots map[*dataaccess.Player]*dataaccess.Player = make(map[*dataaccess.Player]*dataaccess.Player)
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSource)
	tries := 0
	ok := false
	for !ok && tries < 100 {
		remaining := make([]*dataaccess.Player, len(players))
		copy(remaining, players)
		for _, player := range players {
			randomNumber := rnd.Intn(len(remaining))
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
		gamemanagement.saveLots(lots)
		game.Status = dataaccess.Drawn.String()
		gamemanagement.gameRepository.UpdateGame(&game)
	} else {
		log.Warn("Keine plausible Auslosung gefunden")
		drawGameResponseTo.Message = "Nach 100 Versuchen wurde kein plausibles Ergebnis gefunden. Bitte nochmal versuchen oder weniger Ausnahmen definieren."
	}
	drawGameResponseTo.Ok = ok
	return drawGameResponseTo
}

// ResetGame resets a game
func (gamemanagement *gamemanagement) ResetGame(gameCode string) {
	game, _ := gamemanagement.gameRepository.FindGameByCode(gameCode)
	game.Status = dataaccess.Ready.String()
	gamemanagement.gameRepository.UpdateGame(&game)
}

func (gamemanagement *gamemanagement) saveLots(lots map[*dataaccess.Player]*dataaccess.Player) {

	for giftee, gifted := range lots {

		giftee.Gifted = gifted
		gamemanagement.playerRepository.UpdatePlayer(giftee)
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

func (gamemanagement *gamemanagement) generateLink(code string) string {
	return "/game/" + code
}

func (gamemanagement *gamemanagement) generatePassword(plainPassword string) string {
	plainPasswordByte := []byte(plainPassword)
	hash, err := bcrypt.GenerateFromPassword(plainPasswordByte, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Fehler beim Hashen", err)
	}
	return string(hash)
}

func (gamemanagement *gamemanagement) refreshGameStatus(game *dataaccess.Game) {
	_, error := gamemanagement.playerRepository.FindFirstUnreadyPlayerByGameID(game.ID)
	if error != nil {
		game.Status = dataaccess.Ready.String()
		gamemanagement.gameRepository.UpdateGame(game)
	}
}
