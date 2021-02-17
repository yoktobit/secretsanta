package to

// RegisterLoginPlayerPasswordTo ist zum Registrieren eines Spielers
type RegisterLoginPlayerPasswordTo struct {
	GameCode string `json:"gameCode" validate:"required"`
	Name     string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
