package models

type PileReturn struct {
	Discard map[string]int `json:"discard"`
}

type returnCardResponse struct {
	Success   bool       `json:"success"`
	DeckId    string     `json:"deck_id"`
	Remaining int        `json:"remaining"`
	Piles     PileReturn `json:"piles"`
}
