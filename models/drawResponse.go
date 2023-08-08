package models

type Image struct {
	Svg string `json:"svg"`
	Png string `json:"png"`
}

type Card struct {
	Code   string `json:"code"`
	Image  string `json:"image"`
	Images Image  `json:"images"`
	Value  string `json:"value"`
	Suit   string `json:"suit"`
}

type DrawResponse struct {
	Success   bool   `json:"success"`
	DeckId    string `json:"deck_id"`
	Cards     []Card `json:"cards"`
	Remaining int    `json:"remaining"`
}
