package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"log"
	"net/http"
	"os"
)

func init(){
	err := os.Setenv("port", "8080")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := router.GetRouter()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("port"), r))
}