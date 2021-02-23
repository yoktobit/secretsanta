//+build wireinject

package main

import (
	"github.com/google/wire"
	dataaccess "github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/service"
	dataaccess_general "github.com/yoktobit/secretsanta/internal/general/dataaccess"
	logic_general "github.com/yoktobit/secretsanta/internal/general/logic"
)

// InitializeEvent wires together the dependencies
func InitializeEvent() service.RestService {
	wire.Build(service.NewRestService, logic.NewGamemanagement, dataaccess.NewPlayerExceptionRepository, dataaccess.NewPlayerRepository, dataaccess.NewGameRepository, dataaccess_general.NewConnectionWithEnvironment, logic_general.NewRandomizer)
	return service.NewRestService(nil)
}
