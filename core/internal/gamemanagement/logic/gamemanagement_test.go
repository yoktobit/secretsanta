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
	gl "github.com/yoktobit/secretsanta/internal/general/logic"
	"golang.org/x/crypto/bcrypt"
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
		gamemanagement = logic.NewGamemanagement(c, da.NewGameRepository(c), da.NewPlayerRepository(c), da.NewPlayerExceptionRepository(c), gl.NewMockRandomizer())
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
			Expect(createGameResponse.Code).To(MatchRegexp("[a-zA-Z0-9]{21,22}"))
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
		It("should reset a game", func() {
			code := "ABC"
			title := "GameTitle"
			description := "GameDescription"
			status := "Drawn"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "code", "status"}).AddRow(1, title, description, code, status))
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, title, description, code, da.StatusReady.String(), 1).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			err := gamemanagement.ResetGame(code)
			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
			Expect(err).ToNot(HaveOccurred())
		})
		It("should fail to reset a game with empty code", func() {
			code := ""
			err := gamemanagement.ResetGame(code)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Code must not be empty"))
		})
		It("should fail to reset a game if game is not found", func() {
			code := "NonExistingCode"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "code", "status"}))
			err := gamemanagement.ResetGame(code)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
		})
		It("should draw a game randomly", func() {
			code := "ABC"
			title := "GameTitle"
			description := "GameDescription"
			status := "Drawn"
			drawGameTo := to.DrawGameTo{GameCode: code}
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "code", "status"}).AddRow(1, title, description, code, status))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "game_id", "status"}).AddRow(1, "Max", 1, "Ready").AddRow(2, "Moritz", 1, "Ready").AddRow(3, "Susi", 1, "Ready").AddRow(4, "Strolch", 1, "Ready"))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "player_a_id", "player_b_id"}).AddRow(1, 1, 2))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Max"))
			mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "Moritz"))
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, sqlmock.AnyArg(), "", 1, "Ready", "", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, sqlmock.AnyArg(), "", 1, "Ready", "", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, sqlmock.AnyArg(), "", 1, "Ready", "", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, sqlmock.AnyArg(), "", 1, "Ready", "", sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, title, description, code, "Drawn", 1).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			drawGameResponseTo, err := gamemanagement.DrawGame(drawGameTo)
			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
			Expect(err).ToNot(HaveOccurred())
			Expect(drawGameResponseTo.Ok).To(BeTrue())
			Expect(drawGameResponseTo.Message).To(BeEmpty())
		})
		It("should fail to draw a game with too many exceptions", func() {
			code := "ABC"
			title := "GameTitle"
			description := "GameDescription"
			status := "Drawn"
			drawGameTo := to.DrawGameTo{GameCode: code}
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "description", "code", "status"}).AddRow(1, title, description, code, status))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Max").AddRow(2, "Moritz"))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "player_a_id", "player_b_id"}).AddRow(1, 1, 2))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Max"))
			mock.ExpectQuery("SELECT").WithArgs(2).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "Moritz"))
			drawGameResponseTo, err := gamemanagement.DrawGame(drawGameTo)
			Expect(err).ToNot(HaveOccurred())
			Expect(drawGameResponseTo.Ok).To(BeFalse())
			Expect(drawGameResponseTo.Message).To(BeIdenticalTo("Nach 100 Versuchen wurde kein plausibles Ergebnis gefunden. Bitte nochmal versuchen oder weniger Ausnahmen definieren."))
		})
		It("should fail to draw a game with empty code", func() {
			code := ""
			drawGameTo := to.DrawGameTo{GameCode: code}
			drawGameResponseTo, err := gamemanagement.DrawGame(drawGameTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("GameCode must not be empty"))
			Expect(drawGameResponseTo.Ok).To(BeFalse())
			Expect(drawGameResponseTo.Message).To(BeEmpty())
		})
		It("should fail to draw a game if game is not found", func() {
			code := "NonExistentGame"
			drawGameTo := to.DrawGameTo{GameCode: code}
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			drawGameResponseTo, err := gamemanagement.DrawGame(drawGameTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
			Expect(drawGameResponseTo.Ok).To(BeFalse())
			Expect(drawGameResponseTo.Message).To(BeEmpty())
		})
		It("should fail to draw a game if players can't be received", func() {
			code := "ABC"
			drawGameTo := to.DrawGameTo{GameCode: code}
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(gorm.ErrInvalidData)
			drawGameResponseTo, err := gamemanagement.DrawGame(drawGameTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrInvalidData))
			Expect(drawGameResponseTo.Ok).To(BeFalse())
			Expect(drawGameResponseTo.Message).To(BeEmpty())
		})
		It("should fail to draw a game if exceptions can't be received", func() {
			code := "ABC"
			drawGameTo := to.DrawGameTo{GameCode: code}
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Max").AddRow(2, "Moritz"))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(gorm.ErrInvalidData)
			drawGameResponseTo, err := gamemanagement.DrawGame(drawGameTo)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrInvalidData))
			Expect(drawGameResponseTo.Ok).To(BeFalse())
			Expect(drawGameResponseTo.Message).To(BeEmpty())
		})
	})
	Context("Player", func() {
		It("can be added to a game with valid information", func() {
			addRemovePlayerTo := to.AddRemovePlayerTo{Name: "Max", GameCode: "ABC"}
			mock.ExpectQuery("SELECT").WithArgs(addRemovePlayerTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectBegin()
			mock.ExpectQuery("INSERT").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Max", "", 1, "", "Player", nil).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "", "", "", "Waiting", 1).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectCommit()
			err := gamemanagement.AddPlayerToGame(addRemovePlayerTo)
			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
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
			mock.ExpectQuery("SELECT").WithArgs(registerLoginPlayerPasswordTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "Waiting"))
			mock.ExpectQuery("SELECT").WithArgs(registerLoginPlayerPasswordTo.Name, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "role", "game_id"}).AddRow(1, "Max", "", "Player", 1))
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Max", sqlmock.AnyArg(), 1, "Ready", "Player", nil, 1).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("SELECT").WithArgs(1, "Ready").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "", "", "", "Ready", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			err := gamemanagement.RegisterPlayerPassword(registerLoginPlayerPasswordTo)
			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
			Expect(err).ToNot(HaveOccurred())
		})
		It("should fail to register itself if game update fails", func() {
			registerLoginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{Name: "Max", Password: "Pass", GameCode: "ABC"}
			mock.ExpectQuery("SELECT").WithArgs(registerLoginPlayerPasswordTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "Waiting"))
			mock.ExpectQuery("SELECT").WithArgs(registerLoginPlayerPasswordTo.Name, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "role", "game_id"}).AddRow(1, "Max", "", "Player", 1))
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Max", sqlmock.AnyArg(), 1, "Ready", "Player", nil, 1).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("SELECT").WithArgs(1, "Ready").WillReturnError(gorm.ErrInvalidData)
			mock.ExpectCommit()
			err := gamemanagement.RegisterPlayerPassword(registerLoginPlayerPasswordTo)
			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
			Expect(err).ToNot(HaveOccurred())
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
		It("should load the whole list of players by code", func() {
			code := "ABC"
			player1 := to.PlayerResponseTo{Name: "Max", Status: ""}
			player2 := to.PlayerResponseTo{Name: "Moritz", Status: ""}
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Max").AddRow(2, "Moritz"))
			players, err := gamemanagement.GetPlayersByCode(code)
			Expect(err).ToNot(HaveOccurred())
			Expect(players).To(HaveLen(2))
			Expect(players).To(ConsistOf(player1, player2))
		})
		It("should load an empty list of players by code", func() {
			code := "ABC"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))
			players, err := gamemanagement.GetPlayersByCode(code)
			Expect(err).ToNot(HaveOccurred())
			Expect(players).To(BeEmpty())
		})
		It("should not load the list of players by code with empty code", func() {
			code := ""
			players, err := gamemanagement.GetPlayersByCode(code)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Code must not be empty"))
			Expect(players).To(HaveLen(0))
		})
		It("should not load the list of players by code when game not found", func() {
			code := "ABC"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			players, err := gamemanagement.GetPlayersByCode(code)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
			Expect(players).To(HaveLen(0))
		})
		It("should not load the list of players by code when error occured", func() {
			code := "ABC"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(gorm.ErrInvalidData)
			players, err := gamemanagement.GetPlayersByCode(code)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrInvalidData))
			Expect(players).To(HaveLen(0))
		})
		It("should return the role of a player in a game by Code and Name", func() {
			code := "ABC"
			name := "Max"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(name, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "role"}).AddRow(1, "Max", "Admin"))
			role, err := gamemanagement.GetPlayerRoleByCodeAndName(code, name)
			Expect(err).ToNot(HaveOccurred())
			Expect(role).To(BeIdenticalTo(da.RoleAdmin.String()))
		})
		It("should fail to return the role of a player if Code is empty", func() {
			code := ""
			name := "Max"
			role, err := gamemanagement.GetPlayerRoleByCodeAndName(code, name)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Code must not be empty"))
			Expect(role).To(BeEmpty())
		})
		It("should fail to return the role of a player if name is empty", func() {
			code := "ABC"
			name := ""
			role, err := gamemanagement.GetPlayerRoleByCodeAndName(code, name)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Name must not be empty"))
			Expect(role).To(BeEmpty())
		})
		It("should fail to return the role of a player if game is not found", func() {
			code := "NotExistingGame"
			name := "Max"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			role, err := gamemanagement.GetPlayerRoleByCodeAndName(code, name)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
			Expect(role).To(BeEmpty())
		})
		It("should fail to return the role of a player if player is not found", func() {
			code := "ABC"
			name := "NotExistingPlayer"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(name, 1).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			role, err := gamemanagement.GetPlayerRoleByCodeAndName(code, name)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
			Expect(role).To(BeEmpty())
		})
		It("should login a player with valid credentials", func() {
			loginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{GameCode: "ABC", Name: "Max", Password: "12345"}
			plainPasswordByte := []byte(loginPlayerPasswordTo.Password)
			hash, _ := bcrypt.GenerateFromPassword(plainPasswordByte, bcrypt.DefaultCost)
			expectedLoginPlayerPasswordResponseTo := to.RegisterLoginPlayerPasswordResponseTo{Ok: true}
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.Name, 1).WillReturnRows(sqlmock.NewRows([]string{"name", "password"}).AddRow("Max", hash))
			loginPlayerPasswordResponseTo := gamemanagement.LoginPlayer(loginPlayerPasswordTo)
			Expect(loginPlayerPasswordResponseTo).To(BeIdenticalTo(expectedLoginPlayerPasswordResponseTo))
		})
		It("should register a player if he has no password set", func() {
			loginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{GameCode: "ABC", Name: "Max", Password: "12345"}
			expectedLoginPlayerPasswordResponseTo := to.RegisterLoginPlayerPasswordResponseTo{Ok: true}
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "Waiting"))
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.Name, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "role", "game_id"}).AddRow(1, "Max", "", "Player", 1))
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id", "status"}).AddRow(1, "Waiting"))
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.Name, 1).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "password", "role", "game_id"}).AddRow(1, "Max", "", "Player", 1))
			mock.ExpectBegin()
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "Max", sqlmock.AnyArg(), 1, "Ready", "Player", nil, 1).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("SELECT").WithArgs(1, "Ready").WillReturnRows(sqlmock.NewRows([]string{"id"}))
			mock.ExpectExec("UPDATE").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, "", "", "", "Ready", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			mock.ExpectCommit()
			loginPlayerPasswordResponseTo := gamemanagement.LoginPlayer(loginPlayerPasswordTo)
			Expect(loginPlayerPasswordResponseTo).To(BeIdenticalTo(expectedLoginPlayerPasswordResponseTo))
		})
		It("should not login a player if input is invalid", func() {
			loginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{GameCode: "", Name: "", Password: ""}
			loginPlayerPasswordResponseTo := gamemanagement.LoginPlayer(loginPlayerPasswordTo)
			expectedLoginPlayerPasswordResponseTo := to.RegisterLoginPlayerPasswordResponseTo{Ok: false, Message: "Falsche Game-ID, falscher Nutzername oder falsches Passwort"}
			Expect(loginPlayerPasswordResponseTo).To(BeIdenticalTo(expectedLoginPlayerPasswordResponseTo))
		})
		It("should not login a player if game is not found", func() {
			loginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{GameCode: "NotFound", Name: "Max", Password: "12345"}
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			loginPlayerPasswordResponseTo := gamemanagement.LoginPlayer(loginPlayerPasswordTo)
			expectedLoginPlayerPasswordResponseTo := to.RegisterLoginPlayerPasswordResponseTo{Ok: false, Message: "Falsche Game-ID, falscher Nutzername oder falsches Passwort"}
			Expect(loginPlayerPasswordResponseTo).To(BeIdenticalTo(expectedLoginPlayerPasswordResponseTo))
		})
		It("should not login a player if player is not found", func() {
			loginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{GameCode: "NotFound", Name: "Max", Password: "12345"}
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.Name, 1).WillReturnRows(sqlmock.NewRows([]string{"name", "password"}))
			loginPlayerPasswordResponseTo := gamemanagement.LoginPlayer(loginPlayerPasswordTo)
			expectedLoginPlayerPasswordResponseTo := to.RegisterLoginPlayerPasswordResponseTo{Ok: false, Message: "Falsche Game-ID, falscher Nutzername oder falsches Passwort"}
			Expect(loginPlayerPasswordResponseTo).To(BeIdenticalTo(expectedLoginPlayerPasswordResponseTo))
		})
		It("should not login a player if password is wrong", func() {
			loginPlayerPasswordTo := to.RegisterLoginPlayerPasswordTo{GameCode: "NotFound", Name: "Max", Password: "12345"}
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.GameCode).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(loginPlayerPasswordTo.Name, 1).WillReturnRows(sqlmock.NewRows([]string{"name", "password"}).AddRow("Max", "WrongHash"))
			loginPlayerPasswordResponseTo := gamemanagement.LoginPlayer(loginPlayerPasswordTo)
			expectedLoginPlayerPasswordResponseTo := to.RegisterLoginPlayerPasswordResponseTo{Ok: false, Message: "Falsche Game-ID, falscher Nutzername oder falsches Passwort"}
			Expect(loginPlayerPasswordResponseTo).To(BeIdenticalTo(expectedLoginPlayerPasswordResponseTo))
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
		It("should return all exceptions by game code", func() {
			code := "ABC"
			exceptionA := to.ExceptionResponseTo{NameA: "Max", NameB: "Moritz"}
			exceptionB := to.ExceptionResponseTo{NameA: "Susi", NameB: "Strolch"}
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(sqlmock.NewRows([]string{"id", "player_a_id", "player_b_id"}).AddRow(1, 3, 4).AddRow(2, 5, 6))
			mock.ExpectQuery("SELECT").WithArgs(3, 5).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(3, "Max").AddRow(5, "Susi"))
			mock.ExpectQuery("SELECT").WithArgs(4, 6).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(4, "Moritz").AddRow(6, "Strolch"))
			exceptions, err := gamemanagement.GetExceptionsByCode(code)
			Expect(err).ToNot(HaveOccurred())
			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
			Expect(exceptions).To(HaveLen(2))
			Expect(exceptions).To(ConsistOf(exceptionA, exceptionB))
		})
		It("should fail to return all exceptions by empty game code", func() {
			code := ""
			exceptions, err := gamemanagement.GetExceptionsByCode(code)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("Code must not be empty"))
			Expect(exceptions).To(BeEmpty())
		})
		It("should fail to return all exceptions if game was not found", func() {
			code := "NonExistingGame"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}))
			exceptions, err := gamemanagement.GetExceptionsByCode(code)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrRecordNotFound))
			Expect(exceptions).To(BeEmpty())
		})
		It("should fail to return all exceptions if DB error occured", func() {
			code := "ABC"
			mock.ExpectQuery("SELECT").WithArgs(code).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			mock.ExpectQuery("SELECT").WithArgs(1).WillReturnError(gorm.ErrInvalidData)
			exceptions, err := gamemanagement.GetExceptionsByCode(code)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(gorm.ErrInvalidData))
			Expect(exceptions).To(BeEmpty())
		})
	})
})
