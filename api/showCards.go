package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/constants"
	"main.go/initializers"
	"main.go/models"
	"main.go/tools"
)

func ShowPlayerCards(c *gin.Context) {
	//function recieve two parameter from endpoint for list card for player and table----USER ID AND DECK ID
	//we use Deck ID for recognition game of player
	userid := c.Param("userId")
	deckid := c.Param("deckId")

	// create variable type of structure Game to store data from database
	var game models.Game

	//take information for endpoint : deck which user use and name of hand pile for retrieve information from external api
	result := initializers.DB.Where("user_id = ? AND deck_pile = ?", userid, deckid).Find(&game)

	//check error ErrRecordNotFound
	if result.RowsAffected == 0 {
		tools.CheckError(http.StatusBadRequest, c, gorm.ErrRecordNotFound, "Result not found in DB!")
		return
	}

	//call endpoint for list hand cards with necessary information DECK ID and NAME OF HAND PILE used from variable game
	respHand, err := http.Get(fmt.Sprintf(constants.LIST_PILE_CARDS_URL, game.DeckPile, game.HandPile))

	if err != nil {
		tools.CheckError(http.StatusBadRequest, c, err, "Hand cards are not found!")
		//increasing the number of unsuccessful requests to display cards
		tools.UnsuccessfullyShowedCards.Inc()

		return
	}

	// declare variable in acceptable format
	var drawResponse models.ListCardResponse

	//call function for json parse
	parserror := tools.JsonParse(respHand, c, &drawResponse)

	//check parse error
	if parserror != nil {
		tools.CheckError(http.StatusBadRequest, c, parserror, "Error during parse!")
		return
	}

	//variable to store cards for only requested player.
	var handcardsarray []models.CardList

	//http reponse return both json object(hand1 & hand2), we looks for one
	if drawResponse.Piles.Hand1.Cards == nil {
		handcardsarray = drawResponse.Piles.Hand2.Cards
	} else if drawResponse.Piles.Hand2.Cards == nil {
		handcardsarray = drawResponse.Piles.Hand1.Cards
	}

	//call endpoint for list table cards with necessary information DECK ID and NAME OF TABLE PILE used from variable game
	respDeck, err := http.Get(fmt.Sprintf(constants.LIST_PILE_CARDS_URL, game.DeckPile, game.TablePile))

	//handle if there some error from nttp
	if err != nil {
		tools.CheckError(http.StatusBadRequest, c, err, "Table cards are not found!")
		//increasing the number of unsuccessful requests to display cards
        tools.UnsuccessfullyShowedCards.Inc()
		return
	}

	// declare variable in acceptable format
	var drawResponseDeck models.ListCardResponse
	//call function for json parse
	parseerror := tools.JsonParse(respDeck, c, &drawResponseDeck)

	//check parse error
	if parseerror != nil {
		tools.CheckError(http.StatusBadRequest, c, parseerror, "Error during parse!")
		return
	}

	//return response with needed information
	c.JSON(http.StatusOK, gin.H{"User_hand_cards": handcardsarray, "Cards_from_table": drawResponseDeck.Piles.Table})
	//increasing the number of successful requests to display cards
    tools.SuccessfullyShowedCards.Inc()

}

func NewDeckHandler(c *gin.Context) {
	resp, err := http.Get(constants.NEW_DECK_URL)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messege": err,
		})
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messege": err,
		})
	}

	var newDeckResponse models.NewDeckResponse
	json.Unmarshal(body, &newDeckResponse)
	c.JSON(http.StatusOK, gin.H{
		"response": newDeckResponse})
}
