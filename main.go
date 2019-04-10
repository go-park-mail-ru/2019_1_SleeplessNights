package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"net/http"
	"os"
)

func main() {
	//TODO MAKE SECRET MANAGER ?

	//TODO DELETE DATA CREATOR

	err := database.OpenConnection()
	if err != nil {
		logger.Fatal.Print(err.Error())
	}
	defer func() {
		err := database.CloseConnection()
		if err != nil {
			logger.Fatal.Print(err.Error())
		}
	}()
	logger.Info.Printf("\nSuccessfully connected to database at: %s", "5432") //TODO PORT

	faker.CreateFakeData(10)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	logger.Info.Println("Started listening on", PORT)
	r := router.GetRouter()
	logger.Fatal.Println(http.ListenAndServe(":"+PORT, r))
}
