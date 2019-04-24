package thread

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	//Создание форума
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()



	fmt.Println(string(bodyContent))

	type PostData struct {
		Parent  int64  `json:"parent"`
		Author  string `json:"author"`
		Message string `json:"message"`
	}
	var args []PostData

	err = json.Unmarshal(bodyContent, &args)
	if err != nil {
		fmt.Println("Error while parsing request:", err)
		w.WriteHeader(http.StatusInternalServerError)//Возможно, лучше BadRequest
		return
	}

	parents := make([]int64, len(args))
	authors := make([]string, len(args))
	messages := make([]string, len(args))

	for i, postData := range args {
		parents[i]  = postData.Parent
		authors[i]  = postData.Author
		messages[i] = postData.Message
	}

	code, response := createPosts(mux.Vars(r)["slug_or_id"], parents, authors, messages)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error while marshaling response to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if bytes.Equal(responseJSON, []byte("null")) {
		responseJSON = []byte("[]")
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
