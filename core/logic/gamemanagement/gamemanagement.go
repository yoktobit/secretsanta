package logic

import (
	"math/rand"
	"time"

	"github.com/lithammer/shortuuid"
	"github.com/yoktobit/secretsanta/dataaccess"
	"github.com/yoktobit/secretsanta/logic/to"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"

	log "github.com/sirupsen/logrus"
)

// CreateNewGame creates a new Game
func CreateNewGame(createGameTo to.CreateGameTo) to.CreateGameResponseTo {
	code := generateCode()
	game := dataaccess.Game{Code: code, Title: createGameTo.Title, Description: createGameTo.Description, Status: dataaccess.Created.String()}
	dataaccess.DB.Create(&game)
	hashedPassword := generatePassword(createGameTo.AdminPassword)
	player := dataaccess.Player{Name: createGameTo.AdminUser, Password: hashedPassword, GameID: game.ID, Role: dataaccess.RoleAdmin.String(), Status: dataaccess.Ready.String()}
	dataaccess.DB.Create(&player)
	result := to.CreateGameResponseTo{Code: code, Link: generateLink(code)}
	return result
}

// AddPlayerToGame adds a new player to an existing game
func AddPlayerToGame(addPlayerTo to.AddRemovePlayerTo) {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", addPlayerTo.GameCode)
	player := dataaccess.Player{Name: addPlayerTo.Name, GameID: game.ID, Role: dataaccess.RolePlayer.String()}
	dataaccess.DB.Create(&player)
	game.Status = dataaccess.Waiting.String()
	dataaccess.DB.Save(&game)
}

// RemovePlayerFromGame removes an existing player from an existing game
func RemovePlayerFromGame(removePlayerTo to.AddRemovePlayerTo) {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", removePlayerTo.GameCode)
	var player dataaccess.Player
	var exception dataaccess.PlayerException
	result := dataaccess.DB.First(&player, "name = ? AND game_id = ?", removePlayerTo.Name, game.ID)
	if result.Error == nil {
		dataaccess.DB.Delete(&exception, "player_a_id = ? OR player_b_id = ?", player.ID, player.ID)
		dataaccess.DB.Delete(&player, "name = ? AND game_id = ?", removePlayerTo.Name, game.ID)
	}
	refreshGameStatus(&game)
}

// RegisterPlayerPassword registers the password for a player and tells he/she is ready to go
func RegisterPlayerPassword(registerPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) {
	var player dataaccess.Player
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", registerPlayerPasswordTo.GameCode)
	dataaccess.DB.First(&player, "name = ? AND game_id = ?", registerPlayerPasswordTo.Name, game.ID)
	player.Password = generatePassword(registerPlayerPasswordTo.Password)
	player.Status = dataaccess.Ready.String()
	dataaccess.DB.Save(&player)
	refreshGameStatus(&game)
}

// AddException adds a new exception so that PlayerA doesnt have to gift PlayerB
func AddException(addExceptionTo to.AddExceptionTo) {
	var game dataaccess.Game
	var playerA dataaccess.Player
	var playerB dataaccess.Player
	dataaccess.DB.First(&game, "code = ?", addExceptionTo.GameCode)
	dataaccess.DB.First(&playerA, "name = ? AND game_id = ?", addExceptionTo.NameA, game.ID)
	dataaccess.DB.First(&playerB, "name = ? AND game_id = ?", addExceptionTo.NameB, game.ID)
	var existingException dataaccess.PlayerException
	result := dataaccess.DB.First(&existingException, "player_a_id = ? AND player_b_id = ? AND game_id = ?", playerA, playerB, game.ID)
	if result.Error != nil {
		playerException := dataaccess.PlayerException{PlayerA: playerA, PlayerB: playerB, GameID: game.ID}
		dataaccess.DB.Create(&playerException)
	}
}

// GetBasicGameByCode fetches the game from the DB
func GetBasicGameByCode(code string) to.GetBasicGameResponseTo {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", code)
	gameResponseTo := to.GetBasicGameResponseTo{Title: game.Title, Description: game.Description, Code: game.Code}
	return gameResponseTo
}

// GetFullGameByCode fetches the game from the DB
func GetFullGameByCode(code string, playerName string) to.GetFullGameResponseTo {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", code)
	gameResponseTo := to.GetFullGameResponseTo{Title: game.Title, Description: game.Description, Status: game.Status, Code: game.Code}
	if game.Status == dataaccess.Drawn.String() {
		var player dataaccess.Player
		dataaccess.DB.Preload(clause.Associations).First(&player, "name = ? AND game_id = ?", playerName, game.ID)
		gameResponseTo.Gifted = player.Gifted.Name
	}
	return gameResponseTo
}

// GetPlayersByCode fetches the players of a game from the DB
func GetPlayersByCode(code string) []to.PlayerResponseTo {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", code)
	var players []*dataaccess.Player
	dataaccess.DB.Where("game_id = ?", game.ID).Find(&players)
	playerResponseTos := make([]to.PlayerResponseTo, 0)
	for _, player := range players {
		playerResponseTo := to.PlayerResponseTo{Name: player.Name, Status: player.Status}
		playerResponseTos = append(playerResponseTos, playerResponseTo)
	}
	return playerResponseTos
}

// GetPlayerRoleByCodeAndName fetches the player role in a game from the DB
func GetPlayerRoleByCodeAndName(code string, name string) string {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", code)
	var player dataaccess.Player
	result := dataaccess.DB.Where("game_id = ? AND name = ?", game.ID, name).Find(&player)
	if result.Error == nil {
		log.WithField("role", player.Role).Debug("Rolle gefunden")
		return player.Role
	}
	return ""
}

// GetExceptionsByCode returns the draw exceptions in a game
func GetExceptionsByCode(code string) []to.ExceptionResponseTo {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", code)
	var playerExceptions []*dataaccess.PlayerException
	dataaccess.DB.Where("game_id = ?", game.ID).Preload(clause.Associations).Find(&playerExceptions)
	exceptionResponseTos := make([]to.ExceptionResponseTo, 0)
	for _, playerException := range playerExceptions {
		exceptionResponseTo := to.ExceptionResponseTo{NameA: playerException.PlayerA.Name, NameB: playerException.PlayerB.Name}
		exceptionResponseTos = append(exceptionResponseTos, exceptionResponseTo)
	}
	return exceptionResponseTos
}

// LoginPlayer logs in a player
func LoginPlayer(loginPlayerPasswordTo to.RegisterLoginPlayerPasswordTo) to.RegisterLoginPlayerPasswordResponseTo {
	var player dataaccess.Player
	var game dataaccess.Game
	var loginPlayerPasswordResponseTo to.RegisterLoginPlayerPasswordResponseTo
	result := dataaccess.DB.First(&game, "code = ?", loginPlayerPasswordTo.GameCode)
	if result.Error != nil {
		log.WithField("gameCode", loginPlayerPasswordTo.GameCode).Info("Game not found")
		writeLoginError(&loginPlayerPasswordResponseTo)
		return loginPlayerPasswordResponseTo
	}
	result = dataaccess.DB.First(&player, "name = ? AND game_id = ?", loginPlayerPasswordTo.Name, game.ID)
	if result.Error != nil {
		log.WithFields(log.Fields{"name": loginPlayerPasswordTo.Name, "gameCode": loginPlayerPasswordTo.GameCode}).Info("Player not found")
		writeLoginError(&loginPlayerPasswordResponseTo)
		return loginPlayerPasswordResponseTo
	}
	if player.Password == "" {
		RegisterPlayerPassword(loginPlayerPasswordTo)
	} else {
		error := bcrypt.CompareHashAndPassword([]byte(player.Password), []byte(loginPlayerPasswordTo.Password))
		if error != nil {
			writeLoginError(&loginPlayerPasswordResponseTo)
			return loginPlayerPasswordResponseTo
		}
	}
	loginPlayerPasswordResponseTo.Ok = true
	return loginPlayerPasswordResponseTo
}

// DrawGame draws the lots
func DrawGame(drawGameTo to.DrawGameTo) to.DrawGameResponseTo {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", drawGameTo.GameCode)
	var players []dataaccess.Player
	dataaccess.DB.Where(&dataaccess.Player{GameID: game.ID}).Find(&players)
	var exceptions []dataaccess.PlayerException
	dataaccess.DB.Where("game_id = ?", game.ID).Preload(clause.Associations).Find(&exceptions)
	var lots map[dataaccess.Player]dataaccess.Player = make(map[dataaccess.Player]dataaccess.Player)
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSource)
	tries := 0
	ok := false
	for !ok && tries < 100 {
		remaining := make([]dataaccess.Player, len(players))
		copy(remaining, players)
		for _, player := range players {
			randomNumber := rnd.Intn(len(remaining))
			lots[player] = remaining[randomNumber]
			log.WithFields(log.Fields{"giftee": player.Name, "gifted": lots[player].Name}).Debug("Los")
			remaining = append(remaining[:randomNumber], remaining[randomNumber+1:]...)
			log.WithField("remaining", len(remaining)).Debug("Remaining")
			singleOk := isAllowed(exceptions, player, lots[player])
			if !singleOk {
				break
			}
		}
		ok = checkResult(exceptions, lots)
		tries++
	}
	drawGameResponseTo := to.DrawGameResponseTo{}
	if ok {
		saveLots(lots)
		game.Status = dataaccess.Drawn.String()
		dataaccess.DB.Save(&game)
	} else {
		log.Warn("Keine plausible Auslosung gefunden")
		drawGameResponseTo.Message = "Nach 100 Versuchen wurde kein plausibles Ergebnis gefunden. Bitte nochmal versuchen oder weniger Ausnahmen definieren."
	}
	drawGameResponseTo.Ok = ok
	return drawGameResponseTo
}

// ResetGame resets a game
func ResetGame(gameCode string) {
	var game dataaccess.Game
	dataaccess.DB.First(&game, "code = ?", gameCode)
	game.Status = dataaccess.Ready.String()
	dataaccess.DB.Save(&game)
}

func saveLots(lots map[dataaccess.Player]dataaccess.Player) {

	for giftee, gifted := range lots {

		giftee.Gifted = &gifted
		dataaccess.DB.Save(&giftee)
	}
}

func checkResult(exceptions []dataaccess.PlayerException, lots map[dataaccess.Player]dataaccess.Player) bool {

	for giftee, gifted := range lots {
		if !isAllowed(exceptions, giftee, gifted) {
			return false
		}
	}
	return true
}

func isAllowed(exceptions []dataaccess.PlayerException, giftee dataaccess.Player, gifted dataaccess.Player) bool {

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

func writeLoginError(loginPlayerPasswordResponseTo *to.RegisterLoginPlayerPasswordResponseTo) {
	loginPlayerPasswordResponseTo.Message = "Falsche Game-ID, falscher Nutzername oder falsches Passwort"
	loginPlayerPasswordResponseTo.Ok = false
}

func generateCode() string {
	return shortuuid.New()
}

func generateLink(code string) string {
	return "/game/" + code
}

func generatePassword(plainPassword string) string {
	plainPasswordByte := []byte(plainPassword)
	hash, err := bcrypt.GenerateFromPassword(plainPasswordByte, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Fehler beim Hashen", err)
	}
	return string(hash)
}

func refreshGameStatus(game *dataaccess.Game) {
	var otherPlayer dataaccess.Player
	result := dataaccess.DB.First(&otherPlayer, "game_id = ? AND status != ?", game.ID, dataaccess.Ready.String())
	if result.Error != nil {
		game.Status = dataaccess.Ready.String()
		dataaccess.DB.Save(&game)
	}
}
