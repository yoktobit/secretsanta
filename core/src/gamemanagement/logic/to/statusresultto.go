package to

// StatusResultTo contains information about the session status
type StatusResultTo struct {
	LoggedIn bool   `json:"loggedIn"`
	Name     string `json:"name"`
	Role     string `json:"role"`
	GameCode string `json:"gameCode"`
}
