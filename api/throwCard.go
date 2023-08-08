package api

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
	"main.go/constants"
	"main.go/initializers"
	"main.go/models"
	"main.go/tools"
)

func ThrowCardHandler(c *gin.Context) {
	cardCode := c.Param("cardCode")
	deckId := c.Param("deckId")
	playerPile := c.Param("playerPile")
	
	// create variable type of structure Game
	var game models.Game
	result := initializers.DB.Model(&game).Where("hand_pile = ? AND deck_pile = ?", playerPile, deckId).Find(&game)
	tools.ErrorCheck(result.Error, 400, "Failed to fetch data from DB",c)
	//checking if it is the player's turn to play
	if game.First{
		var exist bool = tools.ExistsInDeck(cardCode) 
		//checking if card exist in deck
		if exist {
			//list of cards in a player hand pile
			var HandCards []models.CardList = tools.ListPileCards(deckId, playerPile,c)
			existInHand := tools.ExistsInPile(cardCode, HandCards)

			//checking if the card is in the player's hand
			if existInHand {
				//adding card to table pile
				http.Get(fmt.Sprintf(constants.ADD_TO_PILE_URL, deckId, "table", cardCode)) 
				existInHand = false
				c.JSON(http.StatusOK, gin.H{
					"message": "The card is thrown on the table", 
					"user_hand_cards": tools.ListPileCards(deckId, playerPile, c), 
					"table_cards": tools.ListPileCards(deckId, "table", c),
				})
				tools.WhoPlaysNext(c, playerPile, deckId)
                //increasing number of thrown cards
                tools.SuccessfullyThrownCards.Inc()

				FinishGame(c, deckId)

			} else {
				c.JSON(http.StatusOK, gin.H{"response": "The selected card is not in your hand."})
				//increase number of unsuccessfully thrown cards
                tools.UnsuccessfullyThrownCards.Inc()
			}

		} else {
			c.JSON(http.StatusForbidden, gin.H{"response": "The selected card does not exist in the deck."})
			//increase number of unsuccessfully thrown cards
			tools.UnsuccessfullyThrownCards.Inc()
	}
	}else{
		c.JSON(http.StatusBadRequest, gin.H{"response": "The opponent play next."})
	}
} 
