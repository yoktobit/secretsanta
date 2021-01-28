package main

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	gamemanagement_dataaccess "github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic"
	rest "github.com/yoktobit/secretsanta/internal/gamemanagement/service"
	general_dataaccess "github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

func main() {
	log.SetLevel(log.InfoLevel)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{os.Getenv("ALLOWED_HOSTS")},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	store := cookie.NewStore([]byte(os.Getenv("COOKIE_SECRET")))
	r.Use(sessions.Sessions("mysession", store))
	connection := general_dataaccess.NewConnection()
	gameRepository := gamemanagement_dataaccess.NewGameRepository(connection)
	gamemanagement := logic.NewGamemanagement(gameRepository)
	group := r.Group("/api")
	rest.NewRestService(gamemanagement).DefineRoutes(group)
	r.Run()
}
