package dataaccess

// Role is the role of a player
type Role int

const (
	// RoleAdmin is for a creator or editor of a game
	RoleAdmin Role = iota
	// RolePlayer is a gamer of the game
	RolePlayer
)

func (role Role) String() string {
	return [...]string{"Admin", "Player"}[role]
}
