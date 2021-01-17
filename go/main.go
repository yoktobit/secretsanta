package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/yoktobit/secretsanta/dataaccess"
	"github.com/yoktobit/secretsanta/service/rest"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	store := cookie.NewStore([]byte("secretsantasecret"))
	r.Use(sessions.Sessions("mysession", store))
	dataaccess.ConnectDataBase()
	group := r.Group("/api")
	rest.DefineRoutes(group)
	r.Run()
}
