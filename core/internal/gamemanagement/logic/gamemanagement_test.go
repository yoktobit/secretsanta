package logic_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	da "github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic/to"
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
)

func expectInsertGame(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`INSERT INTO "games"`).WillReturnRows(mock.NewRows([]string{"1"}))
}

func expectInsertPlayer(mock sqlmock.Sqlmock) {
	mock.ExpectQuery(`INSERT INTO "players"`).WillReturnRows(mock.NewRows([]string{"2"}))
}

func expectDefaultQuery(mock sqlmock.Sqlmock) {
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRowsWithColumnDefinition(sqlmock.NewColumn("id")).FromCSVString("1"))
}

func expectDefaultQueryWithNoResult(mock sqlmock.Sqlmock) {
	mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRowsWithColumnDefinition(sqlmock.NewColumn("id")))
}

func expectInsert(mock sqlmock.Sqlmock) {
	mock.ExpectBegin()
	mock.ExpectQuery("INSERT")
	mock.ExpectCommit()
}

var _ = Describe("Gamemanagement", func() {

	var gamemanagement logic.Gamemanagement
	var mock sqlmock.Sqlmock

	BeforeEach(func() {
		var c dataaccess.Connection
		c, mock = dataaccess.NewMockConnection()
		gamemanagement = logic.NewGamemanagement(c, da.NewGameRepository(c), da.NewPlayerRepository(c), da.NewPlayerExceptionRepository(c))
	})

	Context("Game", func() {
		It("can be created with valid information", func() {

			createGameTo := NewCreateGameTo()

			mock.ExpectBegin()
			expectInsertGame(mock)
			expectInsertPlayer(mock)
			mock.ExpectCommit()

			createGameResponse, err := gamemanagement.CreateNewGame(createGameTo)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(createGameResponse).NotTo(BeZero())
			Expect(createGameResponse.Code).NotTo(BeEmpty())
			Expect(createGameResponse.Code).To(HaveLen(22))
		})
		It("fails to be created with an empty object", func() {

			createGameTo := to.CreateGameTo{}

			createGameResponse, err := gamemanagement.CreateNewGame(createGameTo)

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr.Error()).To(BeIdenticalTo(`Key: 'CreateGameTo.Title' Error:Field validation for 'Title' failed on the 'required' tag
Key: 'CreateGameTo.AdminUser' Error:Field validation for 'AdminUser' failed on the 'required' tag
Key: 'CreateGameTo.AdminPassword' Error:Field validation for 'AdminPassword' failed on the 'required' tag`))
			Expect(createGameResponse).To(BeIdenticalTo(to.CreateGameResponseTo{}))
		})
		It("fails to be created without admin user", func() {

			createGameTo := to.CreateGameTo{Title: "Title"}

			createGameResponse, err := gamemanagement.CreateNewGame(createGameTo)

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr.Error()).To(BeIdenticalTo(`Key: 'CreateGameTo.AdminUser' Error:Field validation for 'AdminUser' failed on the 'required' tag
Key: 'CreateGameTo.AdminPassword' Error:Field validation for 'AdminPassword' failed on the 'required' tag`))
			Expect(createGameResponse).To(BeIdenticalTo(to.CreateGameResponseTo{}))
		})
		It("fails to be created without admin pass", func() {

			createGameTo := to.CreateGameTo{Title: "Title", AdminUser: "AdminUser"}

			createGameResponse, err := gamemanagement.CreateNewGame(createGameTo)

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr.Error()).To(BeIdenticalTo(`Key: 'CreateGameTo.AdminPassword' Error:Field validation for 'AdminPassword' failed on the 'required' tag`))
			Expect(createGameResponse).To(BeIdenticalTo(to.CreateGameResponseTo{}))
		})
	})
	Context("Player", func() {
		It("can be added to a game with valid information", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{Name: "Max", GameCode: "1"}
			expectDefaultQuery(mock)
			expectDefaultQuery(mock)
			expectInsert(mock)
			err := gamemanagement.AddPlayerToGame(addRemovePlayerTo)
			Expect(err).ToNot(HaveOccurred())
		})
		It("failes to be added to a game with empty information", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{}
			err := gamemanagement.AddPlayerToGame(addRemovePlayerTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr.Error()).To(BeIdenticalTo(`Key: 'AddRemovePlayerTo.Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'AddRemovePlayerTo.GameCode' Error:Field validation for 'GameCode' failed on the 'required' tag`))
		})
		It("failes to be added to a game with no game code", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{Name: "Max"}
			err := gamemanagement.AddPlayerToGame(addRemovePlayerTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr.Error()).To(BeIdenticalTo(`Key: 'AddRemovePlayerTo.GameCode' Error:Field validation for 'GameCode' failed on the 'required' tag`))
		})
		It("failes to be added to a game which does not exist", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{Name: "Max", GameCode: "ABC"}
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.AddPlayerToGame(addRemovePlayerTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("record not found"))
		})
	})
})
