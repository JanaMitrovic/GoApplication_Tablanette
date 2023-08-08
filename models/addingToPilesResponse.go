package models

type Pile struct {
	Discard map[string]int `json:"discard"`
}

type AddingToPilesResponse struct {
	Success   bool   `json:"success"`
	DeckId    string `json:"deck_id"`
	Remaining int    `json:"remaining"`
	Piles     Pile   `json:"piles"`
}
