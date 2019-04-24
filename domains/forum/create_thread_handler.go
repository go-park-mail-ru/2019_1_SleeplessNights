package forum

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func CreateThreadHandler(w http.ResponseWriter, r *http.Request) {
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println(string(bodyContent))
	args := struct {
		Title   string    `json:"title"`
		Author  string    `json:"author"`
		Message string    `json:"message"`
		Slug    string    `json:"slug"`
		Created time.Time `json:"created"`
	}{}

	err = json.Unmarshal(bodyContent, &args)
	if err != nil {
		if strings.HasPrefix(err.Error(), `parsing time "{}"`) {
			args.Created = time.Time{}
		} else {
			fmt.Println("Error while parsing request:", err)
			w.WriteHeader(http.StatusInternalServerError)//Возможно, лучше BadRequest
			return
		}
	}

	code, response := createThread(mux.Vars(r)["slug"], args.Title, args.Author, args.Message, args.Slug, args.Created)
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
