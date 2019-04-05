package thread

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	code, response := details(mux.Vars(r)["slug_or_id"])
	responseJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error while marshaling response to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	_, err = w.Write(responseJSON)
	if err != nil {
		fmt.Println("Error while writing response body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	return
}
