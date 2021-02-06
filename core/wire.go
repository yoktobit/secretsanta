//+build wireinject

package main

import (
	"github.com/google/wire"
	dataaccess "github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/service"
	dataaccess_general "github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

// InitializeEvent wires together the dependencies
func InitializeEvent() service.RestService {
	wire.Build(service.NewRestService, logic.NewGamemanagement, dataaccess.NewGameRepository, dataaccess_general.NewConnectionWithEnvironment)
	return service.NewRestService(nil)
}
