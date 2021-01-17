package to

// GetFullGameResponseTo gibt Spielinfos zurück
type GetFullGameResponseTo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Gifted      string `json:"gifted"`
	Code        string `json:"code"`
}
