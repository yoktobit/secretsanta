package dataaccess

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the static Database instance
var DB *gorm.DB

// ConnectDataBase connects the database
func ConnectDataBase() {

	dsn := "user=santa password=santa dbname=secretsanta port=5432"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Game{})
	database.AutoMigrate(&Player{})
	database.AutoMigrate(&PlayerException{})

	DB = database
}
