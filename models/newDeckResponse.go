package models

//parital deck
type NewDeckResponse struct {
	Success   bool   `json:"success"`
	DeckId    string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}
