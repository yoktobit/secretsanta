package dataaccess_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	da "github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	gda "github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

var _ = Describe("Repository", func() {

	var repository da.GameRepository
	var connection gda.Connection

	BeforeEach(func() {
		connection = gda.NewConnectionWithConfig(config)
		repository = da.NewGameRepository(connection)
		da.MigrateDb(connection.Connection())
	})
	AfterEach(func() {
	})

	Context("CRUD", func() {
		gameCode := "123"
		game := da.Game{Title: "Title", Description: "Desc", Code: gameCode, Status: da.StatusCreated.String()}
		It("should create a game", func() {
			repository.CreateGame(connection, &game)
			Expect(game.ID).ShouldNot(BeNil())
		})
		It("should read a game", func() {
			game, err := repository.FindGameByCode(gameCode)
			Expect(game.ID).ShouldNot(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("should update a game", func() {
			game.Description = "New Description"
			repository.UpdateGame(connection, &game)
			gameAfterUpdate, err := repository.FindGameByCode(gameCode)
			game.CreatedAt = gameAfterUpdate.CreatedAt
			game.UpdatedAt = gameAfterUpdate.UpdatedAt
			Expect(err).ShouldNot(HaveOccurred())
			Expect(gameAfterUpdate).Should(BeIdenticalTo(game))
		})
	})

})
