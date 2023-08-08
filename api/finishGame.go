package api

import (
	"github.com/gin-gonic/gin"
	"main.go/models"
	"net/http"
	"fmt"
	"strings"
	"main.go/constants"
	"main.go/initializers"
	"main.go/tools"
)



func returnCardsToDeck(deckId string, c *gin.Context){
	retunToDeckURL := fmt.Sprintf(constants.RETURN_TO_DECK_URL, deckId)
	_, errURL := http.Get(retunToDeckURL)
	tools.ErrorCheck(errURL, 400, "Failed API call-Return to deck",c)
}

func numberOfCards(deckId string, c *gin.Context) string {

	Player1TakenPile := tools.Piles(deckId, "taken1",c)
	remainingPlayer1 := Player1TakenPile.Remaining

	Player2TakenPile := tools.Piles(deckId, "taken2",c)
	remainingPlayer2 := Player2TakenPile.Remaining

	if remainingPlayer1 > remainingPlayer2 {
		return "taken1"
	} else if remainingPlayer1 == remainingPlayer2 {
		return "equal"
	} 
	return "taken2"

}

func FinishGame(c *gin.Context, deckId string){

	//Check if player's hands are empty
	//If one of hands is not empty round can continue
	if(tools.NotEmptyHands(deckId,c)){
		c.JSON(http.StatusOK, gin.H{"message": "You can continue round - next player's on turn"})
		return
	}

	//Check if deck is empty
	//If it isn't empty draw draw 2x6 cards from deck to player's hands
	if(!tools.EmptyDeck(deckId, "hand1",c)){
		tools.CreatePile("6", deckId, "hand1", c)
		tools.CreatePile("6", deckId, "hand2", c)
		return
	}

	//If deck is empty round is over
	//Check if table is empty
	if(!tools.EmptyTable(deckId,c)){
		//Draw cards from table
		cardsList := tools.ListPileCards(deckId, "table",c)
		cards := make([]string, 0)
		for _,card := range cardsList{
			cards = append(cards, card.Code)
		}
		cardsString := strings.Join(cards, ",")
		tools.DrawCardsFromPile(deckId, "table", cardsString,c)

		//Check who collected last
		var game models.Game
		result := initializers.DB.Model(&game).Where("deck_pile = ? AND collected_last = true", deckId).Find(&game)
		tools.ErrorCheck(result.Error, 400, "Failed to fetch DB data",c)
		

		//Add to taken pile of the player who collected last
		tools.AddToPile(deckId, game.CollectedPile, cardsString,c)
		
		//Update points
		Score(deckId, game.CollectedPile, cardsString, false,c)

		//Check who has more cards
		playerMoreCards := numberOfCards(deckId,c)

		//update points
		if(playerMoreCards != "equal"){
			var game models.Game
			err := initializers.DB.Where("deck_pile = ? and  collected_pile = ?", deckId, playerMoreCards).Find(&game).Error
			tools.ErrorCheck(err, 400, "Can't find game",c)

			game.Score += 3

			result := initializers.DB.Model(&game).Where("collected_pile = ? AND deck_pile = ?", game.CollectedPile, deckId).Update("score", game.Score)
			tools.ErrorCheck(result.Error, 400, "Error updating score",c)
		}
		
	}

	//Check points
	//If one of the players passed 100 finish game
	var game models.Game
	var games []models.Game
	result := initializers.DB.Model(&game).Where("deck_pile = ? AND collected_last = true", deckId).Find(&games)
	tools.ErrorCheck(result.Error, 400, "Failed to fetch data from DB",c)

	end := false
	for _,game := range games{
		if(game.Score > 100){
			end = true
			break
		}
	}

	if(end){
		result :=initializers.DB.Model(&game).Where("deck_pile = ?", deckId).Update("game_finished", true)
		tools.ErrorCheck(result.Error, 400, "Failed to finish game",c)
		return
	}

	//If game is not finshed - create new round
	//Move all cards from piles to deck
	returnCardsToDeck(deckId,c)
	tools.CreatePile("6", deckId, "hand1", c)
	tools.CreatePile("6", deckId, "hand2", c)
	tools.CreatePile("4", deckId, "table", c)
}