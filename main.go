package main

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"net/http"
	"os"
)

func main() {
	//TODO MAKE SECRET MANAGER ?
	//TODO WRITE TESTS

	//i := 6
	//y := 3
	//
	//s := i % y
	//
	//if s == 0{
	//	fmt.Print(s)
	//}

	PORT :=os.Getenv("PORT")
	if PORT == ""{
		PORT="8080"
	}
	logger.Info.Println("Started listening on", PORT)
	r := router.GetRouter()
	logger.Fatal.Println(http.ListenAndServe(":"+PORT, r))
}
