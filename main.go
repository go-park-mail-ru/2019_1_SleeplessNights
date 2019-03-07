package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"mod/router"
	)

type Config struct {
	port uint
}

var config Config

func init(){
	file, err := os.Open("./config.json")
	if err != nil {
		log.Fatal("Error: Failed to open config.json")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Error: Failed to decode config.json")
	}
}

func main() {
	r := router.GetRouter()
	log.Fatal(http.ListenAndServe(":" + string(config.port), r))
}