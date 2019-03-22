package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/models/forum"
	"io/ioutil"
	"net/http"
)

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	//Создание форума
	bodyContent, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error while reading body:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	args := struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
		User  string `json:"user"`
	}{}

	err = json.Unmarshal(bodyContent, &args)
	if err != nil {
		fmt.Println("Error while parsing request:", err)
		w.WriteHeader(http.StatusInternalServerError)//Возможно, лучше BadRequest
		return
	}

	code, responseData, err := forum.Create(args.Slug, args.Title, args.User)
	if err != nil {
		fmt.Println("Error returned from model:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responseJSON, err := responseData.MarshalToJSON()
	if err != nil {
		fmt.Println("Error while marshaling response data:", err)
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
}

func ForumSlugCreate(w http.ResponseWriter, r *http.Request) {
	//Создание ветки
}

func ForumSlugDetails(w http.ResponseWriter, r *http.Request) {
	//Получение информации о форуме
}

func ForumSlugThreads(w http.ResponseWriter, r *http.Request) {
	//Список ветвей обсуждения форума
}

func ForumSlugUsers(w http.ResponseWriter, r *http.Request) {
	//Пльзователи данного форума
}
