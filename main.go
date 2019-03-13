package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"log"
	"net/http"
	"os"
)

func main() {
	//TODO MAKE LOGGER ?
	//TODO RENAME ROUTER PACKAGE ?
	//TODO MAKE SECRET MANAGER ?
	//TODO WRITE TESTS
	r := router.GetRouter()

	PORT :=os.Getenv("PORT")
	if PORT == ""{
		PORT="8080"
	}
	fmt.Println("Started listening on ", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, r))
}
