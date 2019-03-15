package main

import (
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
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
