package main

import (
	"net/http"
    "main.go/tools"
	"github.com/gin-gonic/gin"
	"main.go/api"
	"main.go/authorization"
	"main.go/initializers"
	"github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"

)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	prometheus.MustRegister(tools.SuccessfullyThrownCards)
    prometheus.MustRegister(tools.UnsuccessfullyThrownCards)
    prometheus.MustRegister(tools.SuccessfullyTakenCards)
    prometheus.MustRegister(tools.UnsuccessfullyTakenCards)
    prometheus.MustRegister(tools.SuccessfullyStartedGame)
    prometheus.MustRegister(tools.UnsuccessfullyStartedGame)
    prometheus.MustRegister(tools.SuccessfullyShowedCards)
    prometheus.MustRegister(tools.UnsuccessfullyShowedCards)
    prometheus.MustRegister(tools.ParsingErrorCounter)
    prometheus.MustRegister(tools.DatabaseErrorCounter)
    prometheus.MustRegister(tools.ServiceStatus)

	r := gin.Default()

	r.POST("/addPlayer", api.AddPlayerHandler)
	r.GET("/cards", api.NewDeckHandler)
	r.GET("/cards/:userId/:deckId", authorization.CheckAuthTokenUser, api.ShowPlayerCards)
	r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", authorization.CheckAuthTokenUser, api.TakeCardsFromTable)
	// r.GET("/takecardsfromtable/:deckId/:handPile/:takenPile", api.TakeCardsFromTable)
	r.GET("/throwCard/:cardCode/:deckId/:playerPile", authorization.CheckAuthTokenUser, api.ThrowCardHandler)
	// r.GET("/throwCard/:cardCode/:deckId/:playerPile", api.ThrowCardHandler)
	r.GET("/gettoken/:userId/:deckId", authorization.MakeToken)

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
    r.GET("/health", healthCheck)

	r.Run()

}
func healthCheck(c *gin.Context) {
    tools.ServiceStatus.Set(1)
    c.JSON(http.StatusOK, gin.H{"status":"OK"})
    
}
