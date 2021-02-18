package logic_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/validator/v10"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	da "github.com/yoktobit/secretsanta/internal/gamemanagement/dataaccess"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic/errors"
	"github.com/yoktobit/secretsanta/internal/gamemanagement/logic/to"
	"github.com/yoktobit/secretsanta/internal/general/dataaccess"
	"gorm.io/gorm"
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
		It("should be found by code", func() {

			code, title, description := "ABC", "GameTitle", "GameDescription"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"code", "title", "description"}).AddRow(code, title, description))
			expectedGetBasicGameResponseTo := to.GetBasicGameResponseTo{Code: code, Title: title, Description: description}

			getBasicGameResponseTo, err := gamemanagement.GetBasicGameByCode(code)

			Expect(err).ToNot(HaveOccurred())
			Expect(getBasicGameResponseTo).To(BeIdenticalTo(expectedGetBasicGameResponseTo))
		})
		It("should not be found with empty code", func() {

			getBasicGameResponseTo, err := gamemanagement.GetBasicGameByCode("")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Code must not be empty"))
			Expect(getBasicGameResponseTo).To(BeIdenticalTo(to.GetBasicGameResponseTo{}))
		})
		It("should not be found if game does not exist", func() {

			code := "anyCode"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"code", "title", "description"}))

			getBasicGameResponseTo, err := gamemanagement.GetBasicGameByCode(code)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
			Expect(getBasicGameResponseTo).To(BeIdenticalTo(to.GetBasicGameResponseTo{}))
		})
		It("should be found as full game with valid information", func() {

			code, title, description, playerName, gifted := "ABC", "GameTitle", "GameDescription", "Max", ""
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "title", "description", "status"}).AddRow(1, code, title, description, da.StatusCreated.String()))
			//mock.ExpectQuery("SELECT SOMETHING").WithArgs(1, "Max").WillReturnRows(sqlmock.NewRows([]string{"gifted"}).AddRow(gifted))
			expectedFullGameResponseTo := to.GetFullGameResponseTo{Code: code, Title: title, Description: description, Status: da.StatusCreated.String(), Gifted: gifted}

			getFullGameResponseTo, err := gamemanagement.GetFullGameByCode(code, playerName)

			Expect(err).ToNot(HaveOccurred())
			Expect(getFullGameResponseTo).To(BeIdenticalTo(expectedFullGameResponseTo))
		})
		It("should be found as full game with gifted player in drawn game with valid information", func() {

			code, title, description, playerName, gifted := "ABC", "GameTitle", "GameDescription", "Max", "Moritz"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "title", "description", "status"}).AddRow(1, code, title, description, da.StatusDrawn.String()))
			mock.ExpectQuery("SELECT").WithArgs(playerName, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "gifted_id"}).AddRow(2, playerName, 3))
			mock.ExpectQuery("SELECT").WithArgs(3).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(3, gifted))
			expectedFullGameResponseTo := to.GetFullGameResponseTo{Code: code, Title: title, Description: description, Status: da.StatusDrawn.String(), Gifted: gifted}

			getFullGameResponseTo, err := gamemanagement.GetFullGameByCode(code, playerName)

			Expect(err).ToNot(HaveOccurred())
			Expect(getFullGameResponseTo).To(BeIdenticalTo(expectedFullGameResponseTo))
		})
		It("should not be found as full game with empty game code", func() {

			getFullGameResponseTo, err := gamemanagement.GetFullGameByCode("", "Max")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Code must not be empty"))
			Expect(getFullGameResponseTo).To(BeIdenticalTo(to.GetFullGameResponseTo{}))
		})
		It("should not be found as full game with empty playerName", func() {

			getFullGameResponseTo, err := gamemanagement.GetFullGameByCode("ABC", "")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("playerName must not be empty"))
			Expect(getFullGameResponseTo).To(BeIdenticalTo(to.GetFullGameResponseTo{}))
		})
		It("should not be found as full game if game does not exist", func() {

			code := "ABC"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "title", "description", "status"}))
			getFullGameResponseTo, err := gamemanagement.GetFullGameByCode(code, "Max")

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
			Expect(getFullGameResponseTo).To(BeIdenticalTo(to.GetFullGameResponseTo{}))
		})
		It("should not be found as full game if player does not exist", func() {

			code, title, description, playerName := "ABC", "GameTitle", "GameDescription", "Max"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id", "code", "title", "description", "status"}).AddRow(1, code, title, description, da.StatusDrawn.String()))
			mock.ExpectQuery("SELECT").WithArgs(playerName, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "gifted_id"}))

			getFullGameResponseTo, err := gamemanagement.GetFullGameByCode(code, playerName)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Player not found"))
			Expect(getFullGameResponseTo).To(BeIdenticalTo(to.GetFullGameResponseTo{}))
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
		It("can be removed from a game with valid information", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{Name: "Max", GameCode: "1"}
			expectDefaultQuery(mock)
			expectDefaultQuery(mock)
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), 1, 1).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), "Max", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			expectDefaultQuery(mock)
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, nil, nil)
			mock.ExpectCommit()
			err := gamemanagement.RemovePlayerFromGame(addRemovePlayerTo)
			Expect(err).ToNot(HaveOccurred())
		})
		It("failes to be removed from a game with empty information", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{}
			err := gamemanagement.RemovePlayerFromGame(addRemovePlayerTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr.Error()).To(BeIdenticalTo(`Key: 'AddRemovePlayerTo.Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'AddRemovePlayerTo.GameCode' Error:Field validation for 'GameCode' failed on the 'required' tag`))
		})
		It("failes to be removed from a game with no game code", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{Name: "Max"}
			err := gamemanagement.RemovePlayerFromGame(addRemovePlayerTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr.Error()).To(BeIdenticalTo(`Key: 'AddRemovePlayerTo.GameCode' Error:Field validation for 'GameCode' failed on the 'required' tag`))
		})
		It("failes to be removed from a game which does not exist", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{Name: "Max", GameCode: "1"}
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.RemovePlayerFromGame(addRemovePlayerTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("record not found"))
		})
		It("failes to be removed from a game where player does not exist", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{Name: "NonExistee", GameCode: "1"}
			expectDefaultQuery(mock)
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.RemovePlayerFromGame(addRemovePlayerTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("record not found"))
		})
		It("should register itself with new credentials", func() {
			registerLoginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{Name: "Max", Password: "Pass", GameCode: "ABC"}
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.RegisterPlayerPassword(registerLoginPlayerPasswordTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("record not found"))
		})
		It("should fail to register itself with empty credentials", func() {
			registerLoginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{}
			err := gamemanagement.RegisterPlayerPassword(registerLoginPlayerPasswordTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr.Error()).To(BeIdenticalTo(`Key: 'RegisterLoginPlayerPasswordTo.GameCode' Error:Field validation for 'GameCode' failed on the 'required' tag
Key: 'RegisterLoginPlayerPasswordTo.Name' Error:Field validation for 'Name' failed on the 'required' tag
Key: 'RegisterLoginPlayerPasswordTo.Password' Error:Field validation for 'Password' failed on the 'required' tag`))
		})
		It("should fail to register itself because the game does not exist", func() {
			registerLoginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{Name: "Max", Password: "Pass", GameCode: "ABC"}
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.RegisterPlayerPassword(registerLoginPlayerPasswordTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
		})
		It("should fail to register itself because the player does not exist", func() {
			registerLoginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{Name: "Max", Password: "Pass", GameCode: "ABC"}
			expectDefaultQuery(mock)
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.RegisterPlayerPassword(registerLoginPlayerPasswordTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
		})
	})
	Context("PlayerException", func() {
		It("should be able to be added to a game", func() {
			addExceptionTo := to.AddExceptionTo{NameA: "Max", NameB: "Erika", GameCode: "ABC"}
			expectDefaultQuery(mock)
			expectDefaultQuery(mock)
			expectDefaultQuery(mock)
			expectDefaultQueryWithNoResult(mock)
			mock.ExpectBegin()
			mock.ExpectQuery("INSERT").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "", "", 0, "", "", nil, 1).WillReturnRows(mock.NewRows([]string{"id"}).AddRow((1)))
			mock.ExpectQuery("INSERT").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "", "", 0, "", "", nil, 1).WillReturnRows(mock.NewRows([]string{"id"}).AddRow((1)))
			mock.ExpectQuery("INSERT").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, 1, 1, 1).WillReturnRows(mock.NewRows([]string{"id"}).AddRow((1)))
			mock.ExpectCommit()
			err := gamemanagement.AddException(addExceptionTo)
			Expect(err).ToNot(HaveOccurred())
		})
		It("should fail to be added with empty information", func() {
			addExceptionTo := to.AddExceptionTo{}
			err := gamemanagement.AddException(addExceptionTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(BeAssignableToTypeOf(validator.ValidationErrors{}))
			valErr := err.(validator.ValidationErrors)
			Expect(valErr).To(HaveLen(3))
		})
		It("should fail to be added because game does not exist", func() {
			addExceptionTo := to.AddExceptionTo{NameA: "Max", NameB: "Erika", GameCode: "ABC"}
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.AddException(addExceptionTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
		})
		It("should fail to be added because playerA does not exist", func() {
			addExceptionTo := to.AddExceptionTo{NameA: "Max", NameB: "Erika", GameCode: "ABC"}
			expectDefaultQuery(mock)
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.AddException(addExceptionTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
		})
		It("should fail to be added because playerB does not exist", func() {
			addExceptionTo := to.AddExceptionTo{NameA: "Max", NameB: "Erika", GameCode: "ABC"}
			expectDefaultQuery(mock)
			expectDefaultQuery(mock)
			expectDefaultQueryWithNoResult(mock)
			err := gamemanagement.AddException(addExceptionTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
		})
		It("should fail to be added because PlayerException already exists", func() {
			addExceptionTo := to.AddExceptionTo{NameA: "Max", NameB: "Erika", GameCode: "ABC"}
			expectDefaultQuery(mock)
			expectDefaultQuery(mock)
			expectDefaultQuery(mock)
			expectDefaultQuery(mock)
			err := gamemanagement.AddException(addExceptionTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(errors.ErrPlayerExceptionAlreadyExists))
		})
	})
})
