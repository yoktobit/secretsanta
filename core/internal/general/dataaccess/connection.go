package dataaccess // import github.com/yoktobit/secretsanta/internal/general/dataacces

import (
	"os"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the static Database instance
var DB *gorm.DB

// ConnectDataBase connects the database
func ConnectDataBase() {

	dsn := os.Getenv("PGSQL_CS")
	log.Infoln("Connecting to " + dsn)
	if dsn == "" {
		dsn = "user=santa password=santa dbname=secretsanta port=5432"
	}
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}
