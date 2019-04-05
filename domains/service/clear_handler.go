package service

import (
	"fmt"
	"net/http"
)

func ClearHandler(w http.ResponseWriter, r *http.Request) {
	err := clear()
	if err != nil {
		fmt.Println("Error while marshaling response to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		return
	}
}
