package dataaccess

import "gorm.io/gorm"

// MigrateDb migrates the DB to the current schema version
func MigrateDb(database *gorm.DB) {

	database.AutoMigrate(&Game{})
	database.AutoMigrate(&Player{})
	database.AutoMigrate(&PlayerException{})
}
