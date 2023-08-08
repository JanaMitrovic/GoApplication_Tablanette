package api

import (
	"strings"
	"github.com/gin-gonic/gin"
	"main.go/initializers"
	"main.go/models"
	"main.go/tools"
)

func Score(deckId string, takenPile string, cards string, table bool, c *gin.Context){
	var game models.Game
	err := initializers.DB.Where("deck_pile = ? and  collected_pile = ?", deckId, takenPile).Find(&game).Error
	tools.ErrorCheck(err, 400,"Error during connecting to base",c)

	var oldScore = game.Score
	var newScore = 0

	//calculating points for every taken card
	codes := strings.Split(cards, ",")
	for _, code := range codes {
		newScore += tools.CalculateCardPoints(code)
	}

	if(table && tools.EmptyTable(deckId,c)){
		newScore++
	}

	game.Score = oldScore + newScore

	result := initializers.DB.Model(&game).Where("collected_pile = ? AND deck_pile = ?", takenPile, deckId).Update("score", game.Score)
	tools.ErrorCheck(result.Error, 400,"Cannot update score",c)
	
}