package forum

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	//Создание форума
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println(string(bodyContent))
	args := struct {
		Slug         string `json:"slug"`
		Title        string `json:"title"`
		UserNickname string `json:"user"`
	}{}

	err = json.Unmarshal(bodyContent, &args)
	if err != nil {
		fmt.Println("Error while parsing request:", err)
		w.WriteHeader(http.StatusInternalServerError)//Возможно, лучше BadRequest
		return
	}

	code, response := create(args.Title, args.UserNickname, args.Slug)
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
