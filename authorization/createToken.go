package authorization

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"main.go/initializers"
	"main.go/models"
	"main.go/tools"
)

func MakeToken(c *gin.Context) {
	//extract parameter from request
	userid := c.Param("userId")
	deckid := c.Param("deckId")

	//find record from database for given parameter from request
	var game models.Game
	result := initializers.DB.Where("user_id = ? AND deck_pile = ?", userid, deckid).Find(&game)

	//check error ErrRecordNotFound
	if result.RowsAffected == 0 {
		tools.CheckError(http.StatusBadRequest, c, gorm.ErrRecordNotFound, "Result not found in DB!")
		return
	}

	//create jwt token with defined payload (user_id, deck_id, exp)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": strconv.Itoa(game.UserID),
		"deck_id": game.DeckPile,
		"exp":     time.Now().Add(time.Hour * 24 * 14).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		tools.CheckError(http.StatusUnauthorized, c, err, "Error occured while createing token")
		return
	}

	//return response with token
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
