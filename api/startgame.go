package api

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
	"main.go/constants"
	"main.go/initializers"
	"main.go/models"
	"main.go/tools"
)

// Making the queue where player will be stored until game starts
var queue []models.User

// Function- adds player to queue.
// Creates a new user according to model, checks for possible errors, adds user to database and to queue
// If enough players are in queue, remove those player from queue and calls function startGame
func AddPlayerHandler(c *gin.Context) {
	var newUser models.User

	err := c.BindJSON(&newUser)
	tools.ErrorCheck(err, 400, "Invalid user data", c)

	result := initializers.DB.Create(&newUser)
	tools.ErrorCheck(result.Error, http.StatusBadRequest, "Failed to create user", c)

	queue = append(queue, newUser)
	if len(queue) >= 2 {
		startGame(queue[0], queue[1], c)
		queue = queue[:2]
	} else {
		c.JSON(201, gin.H{"message": "Wainting for the other player to join"})
	}
}

// Called by addPlayerHandler
// If enough players are present, starts the game
// Creates a game, deck and piles (player hands and table cards)
func startGame(player1 models.User, player2 models.User, c *gin.Context) {
	// Deck creation/alocation
	response, err := http.Get(constants.NEW_SHUFFLED_DECK)
	tools.ErrorCheck(err, http.StatusBadRequest, "Error starting the game-failed to fetch deck", c)
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var deckResponse struct {
			DeckID string `json:"deck_id"`
		}

		err = json.NewDecoder(response.Body).Decode(&deckResponse)
		tools.ErrorCheck(err, http.StatusBadGateway, "Failed to fetch deckID", c)

		//Game creation
		newGame := models.Game{
			Score:         0,
			DeckPile:      deckResponse.DeckID,
			TablePile:     "table",
			HandPile:      "hand1",
			CollectedPile: "taken1",
			First:         true,
			CollectedLast: false,
			UserID:        int(player1.ID),
			User:          player1,
		}
		result := initializers.DB.Create(&newGame)
		tools.ErrorCheck(result.Error, 500, "Failed to add first player into DB ", c)

		newGame2 := models.Game{
			Score:         0,
			DeckPile:      deckResponse.DeckID,
			TablePile:     "table",
			HandPile:      "hand2",
			CollectedPile: "taken2",
			First:         false,
			CollectedLast: false,
			UserID:        int(player2.ID),
			User:          player2,
		}
		result2 := initializers.DB.Create(&newGame2)
        tools.ErrorCheck(result2.Error, 500, "Failed to add second player into DB ", c)

		c.JSON(201, gin.H{"message": "Game has started", "game1": newGame, "game2": newGame2})
		//increase number of successfully started games
		tools.SuccessfullyStartedGame.Inc()


		//Taking 6 cards from deck and forming cards for player1 hands
		tools.CreatePile("6", deckResponse.DeckID, newGame.HandPile, c)

		//Taking 6 cards from deck and forming cards for player2 hands
		tools.CreatePile("6", deckResponse.DeckID, newGame2.HandPile, c)

		//Taking 4 cards from deck and forming cards for table
		tools.CreatePile("4", deckResponse.DeckID, newGame.TablePile, c)

	} else {
		c.JSON(500, gin.H{"message": "Error starting game"})
		//increasing the number of unsuccessful attempts to start the game
		tools.UnsuccessfullyStartedGame.Inc()

	}
}
