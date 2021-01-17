package rest

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	logic "github.com/yoktobit/secretsanta/logic/gamemanagement"
	"github.com/yoktobit/secretsanta/logic/to"
)

// DefineRoutes defines the routes
func DefineRoutes(r *gin.RouterGroup) {

	r.POST("/createNewGame", func(c *gin.Context) {
		session := sessions.Default(c)
		var createGameTo to.CreateGameTo
		c.BindJSON(&createGameTo)
		createGameResponseTo := logic.CreateNewGame(createGameTo)
		session.Set("gameCode", createGameResponseTo.Code)
		session.Set("player", createGameTo.AdminUser)
		session.Save()
		c.JSON(http.StatusOK, createGameResponseTo)
	})
	r.POST("/addPlayer", func(c *gin.Context) {
		session := sessions.Default(c)
		var addPlayerTo to.AddRemovePlayerTo
		c.BindJSON(&addPlayerTo)
		addPlayerTo.GameCode = session.Get("gameCode").(string)
		logic.AddPlayerToGame(addPlayerTo)
		c.Status(http.StatusOK)
	})
	r.POST("/removePlayer", func(c *gin.Context) {
		session := sessions.Default(c)
		var removePlayerTo to.AddRemovePlayerTo
		c.BindJSON(&removePlayerTo)
		removePlayerTo.GameCode = session.Get("gameCode").(string)
		logic.RemovePlayerFromGame(removePlayerTo)
		c.Status(http.StatusOK)
	})

	r.POST("/registerPlayer", func(c *gin.Context) {
		session := sessions.Default(c)
		var registerPlayerPasswordTo to.RegisterLoginPlayerPasswordTo
		c.BindJSON(&registerPlayerPasswordTo)
		registerPlayerPasswordTo.GameCode = session.Get("gameCode").(string)
		logic.RegisterPlayerPassword(registerPlayerPasswordTo)
		c.Status(http.StatusOK)
	})
	r.POST("/loginPlayer", func(c *gin.Context) {
		session := sessions.Default(c)
		var loginPlayerPasswordTo to.RegisterLoginPlayerPasswordTo
		c.BindJSON(&loginPlayerPasswordTo)
		loginPlayerResponseTo := logic.LoginPlayer(loginPlayerPasswordTo)
		if loginPlayerResponseTo.Ok {
			log.Println("Alles ok")
			log.Println(loginPlayerPasswordTo.GameCode)
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
		logic.AddException(addExceptionTo)
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
		drawGameResponseTo := logic.DrawGame(drawGameTo)
		c.JSON(http.StatusOK, drawGameResponseTo)
	})
	r.GET("/reset", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		logic.ResetGame(gameCode.(string))
		c.Status(http.StatusOK)
	})
	r.GET("/game/:gameCode", func(c *gin.Context) {
		gameCode := c.Param("gameCode")
		gameResultTo := logic.GetBasicGameByCode(gameCode)
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
		gameResultTo := logic.GetFullGameByCode(gameCode.(string), player.(string))
		c.JSON(http.StatusOK, gameResultTo)
	})
	r.GET("/players", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		playerResultTos := logic.GetPlayersByCode(gameCode.(string))
		c.JSON(http.StatusOK, playerResultTos)
	})
	r.GET("/exceptions", func(c *gin.Context) {
		session := sessions.Default(c)
		gameCode := session.Get("gameCode")
		if gameCode == nil {
			c.Status(http.StatusForbidden)
			return
		}
		exceptionResponseTos := logic.GetExceptionsByCode(gameCode.(string))
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
			result.Role = logic.GetPlayerRoleByCodeAndName(gameCode.(string), player.(string))
		}
		c.JSON(http.StatusOK, result)
	})
}
