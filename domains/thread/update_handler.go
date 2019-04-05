package thread

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println(string(bodyContent))
	args := struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	}{}

	err = json.Unmarshal(bodyContent, &args)
	if err != nil {
		fmt.Println("Error while parsing request:", err)
		w.WriteHeader(http.StatusInternalServerError)//Возможно, лучше BadRequest
		return
	}

	code, response := update(mux.Vars(r)["slug_or_id"], args.Title, args.Message)
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