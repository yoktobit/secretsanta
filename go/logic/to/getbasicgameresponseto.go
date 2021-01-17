package to

// GetBasicGameResponseTo gibt Spielinfos zurück
type GetBasicGameResponseTo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Code        string `json:"code"`
}
