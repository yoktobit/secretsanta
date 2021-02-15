package logic_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	da "github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic"
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

var _ = Describe("Gamemanagement", func() {

	var gamemanagement logic.Gamemanagement
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var c dataaccess.Connection
		c, mock = dataaccess.NewMockConnection()
		gamemanagement = logic.NewGamemanagement(c, da.NewGameRepository(c), da.NewPlayerRepository(c), da.NewPlayerExceptionRepository(c))
	})

	Context("on blank state with game information", func() {
		createGameTo := NewCreateGameTo()
		It("can create a new game", func() {

			mock.ExpectBegin()
			mock.ExpectQuery(`INSERT INTO "games"`).WillReturnRows(mock.NewRows([]string{"1"}))
			mock.ExpectQuery(`INSERT INTO "players"`).WillReturnRows(mock.NewRows([]string{"2"}))
			mock.ExpectCommit()

			createGameResponse := gamemanagement.CreateNewGame(createGameTo)
			Expect(createGameResponse).NotTo(BeZero())
			Expect(createGameResponse.Code).NotTo(BeEmpty())
			Expect(createGameResponse.Code).To(HaveLen(22))
		})
	})
})
