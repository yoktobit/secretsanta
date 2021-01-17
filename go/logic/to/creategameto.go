package to

// CreateGameTo is for creating a new game
type CreateGameTo struct {
	Title         string `json:"title"`
	Description   string `json:"description"`
	AdminUser     string `json:"adminUser"`
	AdminPassword string `json:"adminPassword"`
}
