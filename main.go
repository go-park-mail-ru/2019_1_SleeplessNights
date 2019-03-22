package main

import (
	"fmt"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/router"
	"net/http"
	"os"
)

func main() {
	PORT :=os.Getenv("PORT")
	if PORT == ""{
		PORT="5000"
	}
	fmt.Println("Started listening on", PORT)
	r := router.GetRouter()
	fmt.Println(http.ListenAndServe(":"+PORT, r))
}
