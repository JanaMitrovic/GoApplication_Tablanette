package models

type DiscardDraw struct {
	Remaining int `json:"remaining"`
}

type PileDraw struct {
	Discard DiscardDraw `json:"discard"`
}

type CardDraw struct {
	Image string `json:"image"`
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type DrawingFromPilesResponse struct {
	Success   bool       `json:"success"`
	DeckId    string     `json:"deck_id"`
	Remaining int        `json:"remaining"`
	Piles     PileDraw   `json:"piles"`
	Cards     []CardDraw `json:"cards"`
}
