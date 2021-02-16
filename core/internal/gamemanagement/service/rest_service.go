package service

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	log "github.com/sirupsen/logrus"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	logic "github.com/yoktobit/secretsanta/internal/gamemanagement/logic"
	to "github.com/yoktobit/secretsanta/internal/gamemanagement/logic/to"
)

// RestService is for defining the REST interface of the app
type RestService interface {
	DefineRoutes(r *gin.RouterGroup)
}

type restService struct {
	gamemanagement logic.Gamemanagement
}

// NewRestService is the factory method for creating the service
func NewRestService(gamemanagement logic.Gamemanagement) RestService {
	return &restService{gamemanagement: gamemanagement}
}

// DefineRoutes defines the routes
func (restService *restService) DefineRoutes(r *gin.RouterGroup) {

	r.POST("/createNewGame", func(c *gin.Context) {
		log.Infoln("createNewGame")
		session := sessions.Default(c)
		log.Infoln("Session ermittelt")
		var createGameTo to.CreateGameTo
		c.BindJSON(&createGameTo)
		createGameResponseTo, err := restService.gamemanagement.CreateNewGame(createGameTo)
		if err != nil {
			_, ok := err.(validator.ValidationErrors)
			if ok {
				c.Status(http.StatusBadRequest)
			} else {
				c.Status(http.StatusInternalServerError)
			}
		}
		log.Infoln("Spiel erstellt")
		session.Clear()
		session.Set("gameCode", createGameResponseTo.Code)
		session.Set("player", createGameTo.AdminUser)
		session.Save()
		log.Infoln("Session gespeichert")
		c.JSON(http.StatusOK, createGameResponseTo)
	})
	r.POST("/addPlayer", func(c *gin.Context) {
		session := sessions.Default(c)
		var addPlayerTo to.AddRemovePlayerTo
		c.BindJSON(&addPlayerTo)
		addPlayerTo.GameCode = session.Get("gameCode").(string)
		restService.gamemanagement.AddPlayerToGame(addPlayerTo)
		c.Status(http.StatusOK)
	})
	r.POST("/removePlayer", func(c *gin.Context) {
		session := sessions.Default(c)
		var removePlayerTo to.AddRemovePlayerTo
		c.BindJSON(&removePlayerTo)
		removePlayerTo.GameCode = session.Get("gameCode").(string)
		restService.gamemanagement.RemovePlayerFromGame(removePlayerTo)
		c.Status(http.StatusOK)
	})

	r.POST("/registerPlayer", func(c *gin.Context) {
		session := sessions.Default(c)
		var registerPlayerPasswordTo to.RegisterLoginPlayerPasswordTo
		c.BindJSON(&registerPlayerPasswordTo)
		registerPlayerPasswordTo.GameCode = session.Get("gameCode").(string)
		restService.gamemanagement.RegisterPlayerPassword(registerPlayerPasswordTo)
		c.Status(http.StatusOK)
	})
	r.POST("/loginPlayer", func(c *gin.Context) {
		session := sessions.Default(c)
		var loginPlayerPasswordTo to.RegisterLoginPlayerPasswordTo
		c.BindJSON(&loginPlayerPasswordTo)
		loginPlayerResponseTo := restService.gamemanagement.LoginPlayer(loginPlayerPasswordTo)
		if loginPlayerResponseTo.Ok {
			log.Infoln("Alles ok")
			log.Infoln(loginPlayerPasswordTo.GameCode)
			session.Set("gameCode", loginPlayerPasswordTo.GameCode)
			session.Set("player", loginPlayerPasswordTo.Name)
		} else {
			session.Clear()
		}
		session.Save()
		c.JSON(http.StatusOK, loginPlayerResponseTo)
	})
	r.POST("/addException", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		var addExceptionTo to.AddExceptionTo
		c.BindJSON(&addExceptionTo)
		addExceptionTo.GameCode = gameCode.(string)
		restService.gamemanagement.AddException(addExceptionTo)
		c.Status(http.StatusOK)
	})
	r.GET("/draw", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		drawGameTo := to.DrawGameTo{GameCode: gameCode.(string)}
		drawGameResponseTo := restService.gamemanagement.DrawGame(drawGameTo)
		c.JSON(http.StatusOK, drawGameResponseTo)
	})
	r.GET("/reset", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		restService.gamemanagement.ResetGame(gameCode.(string))
		c.Status(http.StatusOK)
	})
	r.GET("/game/:gameCode", func(c *gin.Context) {
		gameCode := c.Param("gameCode")
		gameResultTo := restService.gamemanagement.GetBasicGameByCode(gameCode)
		c.JSON(http.StatusOK, gameResultTo)
	})
	r.GET("/game", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		player := session.Get("player")
		if player == nil {
			c.Status(http.StatusForbidden)
			return
		}
		log.Println(gameCode)
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		gameResultTo := restService.gamemanagement.GetFullGameByCode(gameCode.(string), player.(string))
		c.JSON(http.StatusOK, gameResultTo)
	})
	r.GET("/players", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		playerResultTos := restService.gamemanagement.GetPlayersByCode(gameCode.(string))
		c.JSON(http.StatusOK, playerResultTos)
	})
	r.GET("/exceptions", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		exceptionResponseTos := restService.gamemanagement.GetExceptionsByCode(gameCode.(string))
		c.JSON(http.StatusOK, exceptionResponseTos)
	})
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("player", "") // this will mark the session as "written" and hopefully remove the username
		session.Clear()
		session.Save()
		c.Status(http.StatusOK)
	})
	r.GET("/status", func(c *gin.Context) {
		session := sessions.Default(c)
		result := to.StatusResultTo{}
		player := session.Get("player")
		gameCode := session.Get("gameCode")
		if player != nil {
			result.LoggedIn = true
			result.Name = player.(string)
		}
		if gameCode != nil {
			result.GameCode = gameCode.(string)
		}
		if player != nil && gameCode != nil {
			result.Role = restService.gamemanagement.GetPlayerRoleByCodeAndName(gameCode.(string), player.(string))
		}
		c.JSON(http.StatusOK, result)
	})
}
