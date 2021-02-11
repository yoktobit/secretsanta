// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	dataaccess2 "github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/service"
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

// Injectors from wire.go:

// InitializeEvent wires together the dependencies
func InitializeEvent() service.RestService {
	connection := dataaccess.NewConnectionWithEnvironment()
	gameRepository := dataaccess2.NewGameRepository(connection)
	playerRepository := dataaccess2.NewPlayerRepository(connection)
	playerExceptionRepository := dataaccess2.NewPlayerExceptionRepository(connection)
	gamemanagement := logic.NewGamemanagement(gameRepository, playerRepository, playerExceptionRepository)
	restService := service.NewRestService(gamemanagement)
	return restService
}
