package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"log"
	"net/http"
	"os"
)

func main() {
	r := router.GetRouter()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), r))
}
