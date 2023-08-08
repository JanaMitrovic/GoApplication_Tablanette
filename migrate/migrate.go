package main

import (
	"main.go/initializers"
	"main.go/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()

}

func main() {
	initializers.DB.AutoMigrate(&models.User{}, &models.Game{})
}
