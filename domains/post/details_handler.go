package post

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

func DetailsHandler(w http.ResponseWriter, r *http.Request) {
	related := strings.Split(r.URL.Query().Get("related"), ",")
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		fmt.Println("Error while parsing mux variables:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	code, response := details(int64(id), related)
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