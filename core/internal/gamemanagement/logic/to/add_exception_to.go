package to

// AddExceptionTo gibt eine Ausnahme im Spiel an
type AddExceptionTo struct {
	NameA    string `json:"nameA" validate:"required"`
	NameB    string `json:"nameB" validate:"required"`
	GameCode string `json:"gameCode" validate:"required"`
}
