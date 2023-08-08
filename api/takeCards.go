package api

import (
	"github.com/gin-gonic/gin"
	"main.go/models"
	"net/http"
	"strings"
	"strconv"
	"main.go/initializers"
	"main.go/tools"
)

//Function that changes who collected last
func changeWhoCollectedLast(c *gin.Context, handPile string, deckId string){
	var game models.Game
	//set attribute "collected_last" on false for player 1
	result :=initializers.DB.Model(&game).Where("hand_pile = ? AND deck_pile = ?", handPile, deckId).Update("collected_last", true)
	tools.ErrorCheck(result.Error,400,"Failed DB update",c)

	//set attribute "collected_last" on true for player 2
	result =initializers.DB.Model(&game).Where("hand_pile NOT IN (?) AND deck_pile = ?", handPile, deckId).Update("collected_last", false)
	tools.ErrorCheck(result.Error,400,"Failed DB update",c)
}

type RequestData struct{
	HandCard string `json:"hand_card"`
	TakenCards string `json:"taken_cards"`
}

//Function for checking if sum of cards group is same as hand card value
func isGroupValid(c *gin.Context, HandCard string, TakenCards []string)(valid bool){
	valid = false

	//CHECK VALUES
	var HandCardValue int
	var sum int = 0
	var countA int = 0
	switch(HandCard[0]){
	case '0': HandCardValue = 10
	case 'A': HandCardValue = 11 //We will always consider A from hand as 11
	case 'J': HandCardValue = 12
	case 'Q': HandCardValue = 13
	case 'K': HandCardValue = 14
	default: HandCardValue,_ = strconv.Atoi(string(HandCard[0]))
	}

	//Count number of As in Taken pile
	for _, cardTaken := range TakenCards{
		if(cardTaken[0] == 'A'){
			countA++;
		}
	}

	//depending on nomber of As in Taken pile we have 5 cases
	//Case1: countA=0
	//Case2: countA=1 -> 2 combinations: 1             sum=1
	//                                   11            sum=11 
	
	//Case3: countA=2 -> 3 combinations: 1  1          sum=2
	//                                   1  11         sum=12 
	//                                   11 11         sum=22

	//Case4: countA=3 -> 4 combinations: 1  1  1       sum=3 
	//                                   1  1  11      sum=13
	//                                   1  11 11      sum=23
	//                                   11 11 11      sum=33

	//Case5: countA=4 -> 5 combinations: 1  1  1  1    sum=4
	//                                   1  1  1  11   sum=14
	//                                   1  1  11 11   sum=24
	//                                   1  11 11 11   sum=34
	//                                   11 11 11 11   sum=44

	for _, cardTaken := range TakenCards{
		var val int
		switch(cardTaken[0]){
			case '0': val = 10
			case 'A': val = 1 //Firstly, we will consider all As as 1 and calculate the sum
			case 'J': val = 12
			case 'Q': val = 13
			case 'K': val = 14
			default: val,_ = strconv.Atoi(string(cardTaken[0]))
		}
		sum += val
	}

	//As the sum of each combination differs by 10, we will check countA+1 sums each greater by 10
	//If modul of any sum by HandCardValue is 0, player can take cards otherwise he can't
	for i:=0; i<=countA; i++{
		sum += i*10;
		if (sum == HandCardValue){
			valid = true
			return
		}
	}

	return
}


func TakeCardsFromTable(c *gin.Context){
	//EXTRACT PARAMETERS
	deckId := c.Param("deckId")
	handPile := c.Param("handPile")
	takenPile := c.Param("takenPile")

	//CHECK IF IT PLAYER'S TURN
	var game models.Game
	result := initializers.DB.Model(&game).Where("hand_pile = ? AND deck_pile = ?", handPile, deckId).Find(&game)
	tools.ErrorCheck(result.Error,400,"Failed to find game",c)

	if !game.First{
		c.JSON(http.StatusBadRequest, gin.H{"response": "The opponent plays next."})
		return
	}

	//EXTRACT BODY REQUEST
	var RequestData RequestData
	err := c.BindJSON(&RequestData)
	tools.ErrorCheck(err,400,"Failed to read body",c)

	HandCard := RequestData.HandCard
	TakenCardsString := RequestData.TakenCards
	TakenCardsGroups := strings.Split(TakenCardsString, ";")
	var TakenCards []string

	var valid bool = true

	//VALIDATE EACH CARDS GROUP
	for _, group := range TakenCardsGroups{
		TakenCards = strings.Split(group, ",")
		
		//VALIDATE CARD FROM HAND
		if(!tools.ExistsInDeck(HandCard)){
			c.JSON(http.StatusForbidden, gin.H{"response": "The selected hand card does not exist in the deck."})
			//increase number of unsuccessfully taken cards
			tools.UnsuccessfullyTakenCards.Inc()
			return
		}

		var HandCards []models.CardList = tools.ListPileCards(deckId, handPile,c)

		if(!tools.ExistsInPile(HandCard, HandCards)){
			c.JSON(http.StatusForbidden, gin.H{"response": "The selected card is not in your hand."})
			//increase number of unsuccessfully taken cards
			tools.UnsuccessfullyTakenCards.Inc()
			return
		}

		//VALIDATE CARDS FROM TABLE
		var TableCards []models.CardList = tools.ListPileCards(deckId, "table",c)
		for _, cardTaken := range TakenCards {
			if(!tools.ExistsInDeck(cardTaken)){
				c.JSON(http.StatusForbidden, gin.H{"response": "The selected table card does not exist in the deck."})
				//increase number of unsuccessfully taken cards
				tools.UnsuccessfullyTakenCards.Inc()
				return
			}
			
			if(!tools.ExistsInPile(cardTaken, TableCards))	{
				c.JSON(http.StatusForbidden, gin.H{"response": "Some of selected cards is not on the table."})
				//increase number of unsuccessfully taken cards
				tools.UnsuccessfullyTakenCards.Inc()
				return
			}	
		}

		if(!isGroupValid(c, HandCard, TakenCards)){
			valid = false
			break
		}

	}

	//IF ONE OF THE CARDS GROUP IS NOT VALID
	if(!valid){
		c.JSON(http.StatusNotFound, gin.H{"response": "You can't take chosen cards"})
		//increase number of unsuccessfully taken cards
		tools.UnsuccessfullyTakenCards.Inc()
		return
	}

	//IF VALID MOVE CARDS FROM HAND AND TABLE PILE TO TAKEN PILE
	tools.DrawCardsFromPile(deckId, handPile, HandCard,c)
	cards := strings.Join(TakenCardsGroups, ",")
	tools.DrawCardsFromPile(deckId, "table", cards,c)
	tools.AddToPile(deckId, takenPile, cards+","+HandCard,c)

	//NOTE THAT THIS PLAYER HAS COLLECTED CARDS LAST AND CHANGE WHO PLAYS NEXT
	tools.WhoPlaysNext(c, handPile, deckId)
	changeWhoCollectedLast(c, handPile, deckId)

	c.JSON(http.StatusOK, gin.H{
		"response": "Cards are moved from hand and table pile to taken pile",
		"user_hand_cards": tools.ListPileCards(deckId, handPile, c), 
		"table_cards": tools.ListPileCards(deckId, "table", c),
	})

	Score(deckId, takenPile, cards + "," + HandCard, true,c)
	FinishGame(c, deckId)
	
	//increase number of successfully taken cards
	tools.SuccessfullyTakenCards.Inc()

}