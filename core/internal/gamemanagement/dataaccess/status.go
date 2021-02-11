package dataaccess

// Status is the status of a player or game
type Status int

const (
	// StatusCreated is for newly created Players/Games
	StatusCreated Status = iota
	// StatusWaiting is only for games, for waiting games for having all players registered
	StatusWaiting
	// StatusReady is for players and games, a game is ready when all players are
	StatusReady
	// StatusDrawn is for games only and describe a final state for games, where all lots are drawn
	StatusDrawn
	// StatusReset is for resetted games, same as Ready, but only possible after a draw
	StatusReset
)

func (status Status) String() string {
	return [...]string{"Created", "Waiting", "Ready", "Drawn", "Reset"}[status]
}
