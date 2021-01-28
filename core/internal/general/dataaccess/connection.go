package dataaccess // import github.com/yoktobit/secretsanta/internal/general/dataacces

import (
	"os"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connection encapsulates some DB connection
type Connection interface {
	Connection() *gorm.DB
}

type connection struct {
	db *gorm.DB
}

// NewConnection connects the database
func NewConnection() Connection {

	dsn := os.Getenv("PGSQL_CS")
	log.Infoln("Connecting to " + dsn)
	if dsn == "" {
		dsn = "user=santa password=santa dbname=secretsanta port=5432"
	}
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	return &connection{db: database}
}

// Connection gets the Gorm-Connection from the Connection object
func (connection *connection) Connection() *gorm.DB {

	return connection.db
}
