package thread

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	fmt.Println(string(bodyContent))

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil {
		limit = 100
	}

	since, err := strconv.Atoi(r.URL.Query().Get("since"))
	if err != nil {
		since = 0
	}

	sort := r.URL.Query().Get("sort")
	if sort == "" {
		sort = "flat"
	}

	desc, err := strconv.ParseBool(r.URL.Query().Get("desc"))
	if err != nil {
		desc = false
	}

	code, response := posts(mux.Vars(r)["slug_or_id"], int32(limit), int64(since), sort, desc)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error while marshaling response to JSON:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if bytes.Equal(responseJSON,[]byte("null")) {
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
