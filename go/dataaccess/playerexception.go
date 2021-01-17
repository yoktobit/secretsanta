package dataaccess

import "gorm.io/gorm"

// PlayerException is an exception so that NameA doesnt have to gift NameB
type PlayerException struct {
	gorm.Model
	PlayerAID uint
	PlayerBID uint
	PlayerA   Player `gorm:"foreignKey:PlayerAID"`
	PlayerB   Player `gorm:"foreignKey:PlayerBID"`
	GameID    uint
}
