package initializers

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
	var host = os.Getenv("HOST")
	var port, _ = strconv.Atoi(os.Getenv("DB_PORT"))
	var user = os.Getenv("USER")
	var password = os.Getenv("PASSWORD")
	var dbname = os.Getenv("DB_NAME")
	var err error

	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database")
	}

}
