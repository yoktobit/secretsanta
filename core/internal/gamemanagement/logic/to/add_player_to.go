package to

// AddRemovePlayerTo is for adding new player to a game
type AddRemovePlayerTo struct {
	Name     string `json:"name" validate:"required"`
	GameCode string `json:"gameCode" validate:"required"`
}
