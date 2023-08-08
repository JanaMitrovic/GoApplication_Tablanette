package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `json:"name"`
}

type Game struct {
	gorm.Model
	Score         int    `json:"score"`
	DeckPile      string `json:"deckPile"`
	TablePile     string `json:"tablePile"`
	HandPile      string `json:"handPile"`
	CollectedPile string `json:"collectedPile"`
	First         bool   `json:"first"`
	CollectedLast bool   `json:"collectedLast"`
	GameFinished  bool   `json:"gameFinished"`
	UserID        int    `json:"user_id"`
	User          User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
