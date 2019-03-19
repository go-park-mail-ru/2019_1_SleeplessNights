package main

import (
	"github.com/DragonF0rm/Technopark-DBMS-Forum/logger"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/router"
	"net/http"
	"os"
)

func main() {
	//TODO MAKE SECRET MANAGER ?
	//TODO WRITE TESTS

	PORT :=os.Getenv("PORT")
	if PORT == ""{
		PORT="5000"
	}
	logger.Info.Println("Started listening on", PORT)
	r := router.GetRouter()
	logger.Fatal.Println(http.ListenAndServe(":"+PORT, r))
}
