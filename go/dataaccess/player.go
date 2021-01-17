package dataaccess

import "gorm.io/gorm"

// Player is a player in the game
type Player struct {
	gorm.Model
	Name     string `json:"name"`
	Password string `json:"password"`
	GameID   uint
	Status   string
	Role     string
	GiftedID *uint
	Gifted   *Player `gorm:"foreignKey:GiftedID"`
}
