package to

// DrawGameResponseTo is the answer sent after drawing lots
type DrawGameResponseTo struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
