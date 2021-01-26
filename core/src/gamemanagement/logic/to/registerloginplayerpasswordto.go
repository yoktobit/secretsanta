package to

// RegisterLoginPlayerPasswordTo ist zum Registrieren eines Spielers
type RegisterLoginPlayerPasswordTo struct {
	GameCode string `json:"gameCode"`
	Name     string `json:"username"`
	Password string `json:"password"`
}
