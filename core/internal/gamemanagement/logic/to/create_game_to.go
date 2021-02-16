package to

// CreateGameTo is for creating a new game
type CreateGameTo struct {
	Title         string `json:"title" validate:"required"`
	Description   string `json:"description"`
	AdminUser     string `json:"adminUser" validate:"required"`
	AdminPassword string `json:"adminPassword" validate:"required"`
}
