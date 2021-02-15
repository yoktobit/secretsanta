package dataaccess

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewMockConnection initializes a new mock connection
func NewMockConnection() (Connection, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic("Sql-Mock could not be initialized")
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), nil)
	return &connection{db: gormDB}, mock
}
