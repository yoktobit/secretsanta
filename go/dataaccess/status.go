package dataaccess

// Status is the status of a player or game
type Status int

const (
	// Created is for newly created Players/Games
	Created Status = iota
	// Waiting is only for games, for waiting games for having all players registered
	Waiting
	// Ready is for players and games, a game is ready when all players are
	Ready
	// Drawn is for games only and describe a final state for games, where all lots are drawn
	Drawn
	// Reset is for resetted games, same as Ready, but only possible after a draw
	Reset
)

func (status Status) String() string {
	return [...]string{"Created", "Waiting", "Ready", "Drawn", "Reset"}[status]
}
