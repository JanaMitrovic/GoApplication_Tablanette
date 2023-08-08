package tools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"main.go/initializers"
	"main.go/constants"
	"main.go/models"
	"github.com/gin-gonic/gin"
)

//Function for parsing response JSON to Struct
func ParseJsonToStruct(resp *http.Response, c *gin.Context)(body []byte){
	defer resp.Body.Close()

	body, errBody := io.ReadAll(resp.Body)
	ErrorCheck(errBody,400,"Failed to parse data",c)

	return
}
// Helper function for error responses
func ErrorCheck(err error, errCode int, errMsg string, c *gin.Context) {
	if err != nil {
		c.JSON(errCode, gin.H{"message": errMsg})
		//increase number of parsing error 
		ParsingErrorCounter.Inc()

		return
	}
}

func AddToPile(deckID string, pileToBeCreated string, cardCodes string, c *gin.Context){
	addToPileURL := fmt.Sprintf(constants.ADD_TO_PILE_URL, deckID, pileToBeCreated, cardCodes)
	newPile, err := http.Get(addToPileURL)
	ErrorCheck(err, 500, "Failed API call-Add to Pile", c)
	defer newPile.Body.Close()

	if newPile.StatusCode == http.StatusOK {
		var player1HandPileResponse models.AddingToPilesResponse

		err = json.NewDecoder(newPile.Body).Decode(&player1HandPileResponse)
		ErrorCheck(err, 400, "Failed to fetch data from API call", c)
	}
}

//Function for listing cards in a pile
func ListPileCards(deck string, pileName string, c *gin.Context)(CardsArray []models.CardList){
	listPileCardsURL := fmt.Sprintf(constants.LIST_PILE_CARDS_URL, deck, pileName)
	resp, errURL := http.Get(listPileCardsURL)
	ErrorCheck(errURL,500,"Faile API call - List pile cards",c)
	defer resp.Body.Close()

	body := ParseJsonToStruct(resp,c) 

	var ListCardResponse models.ListCardResponse
	err := json.Unmarshal(body, &ListCardResponse)
	ErrorCheck(err,400,"Faile to fetch data from API",c)

	switch(pileName){
	 case "hand1": CardsArray = ListCardResponse.Piles.Hand1.Cards
	 case "hand2": CardsArray = ListCardResponse.Piles.Hand2.Cards
	 case "table": CardsArray = ListCardResponse.Piles.Table.Cards
	 default: 
	}

	return
}

//Function for drawing cards from a pile
func DrawCardsFromPile(deck string, pileName string, cards string, c *gin.Context){
	drawCardsFromPileURL := fmt.Sprintf(constants.DRAW_CARDS_FROM_PILE_URL, deck, pileName, cards)
	resp, errURL := http.Get(drawCardsFromPileURL)
	ErrorCheck(errURL,500,"Faile API call - Draw cards from pile",c)
	defer resp.Body.Close()

	body := ParseJsonToStruct(resp,c) 

	var DrowCardResponse models.DrawingFromPilesResponse
	err := json.Unmarshal(body, &DrowCardResponse)
	ErrorCheck(err,400,"Faile to fetch data from API",c)
}


// Used to create 3 piles at the start of game
// Those 3 piles are player hands(for each player) and table pile
func CreatePile(numberOfCards string, deckID string, pileToBeCreated string, c *gin.Context) {
	draw_A_Card := fmt.Sprintf(constants.DRAW_A_CARD_URL, deckID, numberOfCards)
	drawnCards, err := http.Get(draw_A_Card)
	errorMessage := "Error starting the game"
	ErrorCheck(err, 500, errorMessage, c)
	defer drawnCards.Body.Close()

	if drawnCards.StatusCode == http.StatusOK {

		var drawnCardsResponse struct {
			Success   bool          `json:"success"`
			DeckId    string        `json:"deck_id"`
			Cards     []models.Card `json:"cards"`
			Remaining int           `json:"remaining"`
		}

		err = json.NewDecoder(drawnCards.Body).Decode(&drawnCardsResponse)
		ErrorCheck(err, 500, errorMessage, c)

		//Taking card codes for API URL
		cardCodes := ""
		for i := 0; i < len(drawnCardsResponse.Cards); i++ {
			cardCodes += drawnCardsResponse.Cards[i].Code
			if i < len(drawnCardsResponse.Cards)-1 {
				cardCodes += ","
			}
		}
		AddToPile(deckID,pileToBeCreated,cardCodes,c)
	}
}

//Function for checking if card exists in deck
func ExistsInDeck(cardCode string)(exist bool){
	var allCards = [52]string{"AS","2S","3S","4S","5S","6S","7S","8S","9S","0S","JS","QS","KS",
	"AD","2D","3D","4D","5D","6D","7D","8D","9D","0D","JD","QD","KD",
	"AC","2C","3C","4C","5C","6C","7C","8C","9C","0C","JC","QC","KC",
	"AH","2H","3H","4H","5H","6H","7H","8H","9H","0H","JH","QH","KH"}

	exist = false
	for i:=0; i < 52; i++ {
		if allCards[i]==cardCode{
			exist=true
		}
	}

	return
}

//Function for checking if card exists in specific pile
func ExistsInPile(cardCode string, pile []models.CardList)(exist bool){
	exist = false
	for _, card := range pile {
		if card.Code == cardCode{
			exist=true
		}
	}
	return
}

func WhoPlaysNext(c *gin.Context, playerPile string, deckId string){
	var game models.Game
	//set attribute "first" on false for player 1
	result :=initializers.DB.Model(&game).Where("hand_pile = ? AND deck_pile = ?", playerPile, deckId).Update("first", false)
	ErrorCheck(result.Error, 500, "Failed to update DB", c)

	//set attribute "first" on true for player 2
	result =initializers.DB.Model(&game).Where("hand_pile NOT IN (?) AND deck_pile = ?", playerPile, deckId).Update("first", true)
	ErrorCheck( result.Error, 500, "Failed to update DB", c)
}

func Piles(deckId string, pile string, c *gin.Context) models.PileList {

	var listCardResponse = ListCardResponseFunction(deckId, pile,c)
	var Pile models.PileList

	if pile == "taken1" {
		Pile = listCardResponse.Piles.Taken1
	} else if pile == "taken2" {
		Pile = listCardResponse.Piles.Taken2
	} else if pile == "hand1" {
		Pile = listCardResponse.Piles.Hand1
	} else if pile == "hand2" {
		Pile = listCardResponse.Piles.Hand2
	} else if pile == "table" {
		Pile = listCardResponse.Piles.Table
	}
	return Pile
}

func NotEmptyHands(deckId string, c *gin.Context) bool {

	var Pile1 models.PileList = Piles(deckId, "hand1",c)
	remaining1 := Pile1.Remaining

	var Pile2 models.PileList = Piles(deckId, "hand2",c)
	remaining2 := Pile2.Remaining

	if remaining1 != 0 && remaining2 != 0 {
		return true
	}
	return false
}

func ListCardResponseFunction(deckId string, pile string, c *gin.Context) models.ListCardResponse {
	listPileCardsURL := fmt.Sprintf(constants.LIST_PILE_CARDS_URL, deckId, pile)

	resp, err := http.Get(listPileCardsURL)
	ErrorCheck( err, 404, "Issue with specified URL", c)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	ErrorCheck( err, 500, "Error during reading response!", c)
	
	var listCardResponse models.ListCardResponse
	err = json.Unmarshal(body, &listCardResponse)
	ErrorCheck( err, 500, "Failed to read body", c)

	return listCardResponse
}

func EmptyDeck(deckId string, pile string, c *gin.Context) bool {
	var listCardResponse = ListCardResponseFunction(deckId, pile,c)
	remaining := listCardResponse.Remaining

	if remaining == 0 {
		return true
	}
	return false
}

func EmptyTable(deckId string, c *gin.Context) bool {
	TablePile := Piles(deckId, "table", c)
	remainingTable := TablePile.Remaining

	if remainingTable == 0 {
		return true
	}
	return false
}

func CalculateCardPoints(code string) int {

	if code[0] == 'A' || code[0] == 'K' || code[0] == 'Q' || code[0] == 'J' {
		return 1
	} else if code[0] == '2' && code[1] == 'C' {
		return 1
	} else if code[0] == '0' {
		if code[1] == 'D' {
			return 2
		}
		return 1
	}
	return 0
}

func JsonParse(response *http.Response, c *gin.Context, parameter interface{}) error {
	//part of code for accepting response
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	//check error during reading response
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Error during reading response!"})
		//increase number of parsing error 
		ParsingErrorCounter.Inc()

	}

	// error expectation if exist
	errHandUm := json.Unmarshal(body, parameter)

	return errHandUm
}

func CheckError(code int, c *gin.Context, err error, message string) {

	if err != nil {
		if message != "" {
			c.JSON(code, gin.H{"message": message})
		} else {
			c.JSON(code, gin.H{"message": err})
		}
	}
}