package dataaccess

import "gorm.io/gorm"

// Game is the main game
type Game struct {
	gorm.Model
	//ID    uint   `json:"id" gorm:"primary_key"`
	Title       string `json:"title"`
	Description string
	Code        string `json:"author"`
	Status      string `json:"status"`
}
