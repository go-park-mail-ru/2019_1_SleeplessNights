package handlers

import (
	"encoding/json"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/handlers/helpers"
	"github.com/DragonF0rm/Technopark-DBMS-Forum/models"
	"net/http"
	"sort"
	"strconv"
)

const (
	pagesPerList = 4
)

type LeaderBoard struct {
	PagesTotal  int           `json:"pages_total"`
	CurrentPage int           `json:"page"`
	Slice       []models.User `json:"data"`
}

func Paginate(data []interface{}, skip int) []interface{} {
	//TODO MAKE UNIVERSAL AND MOVE TO HELPERS
	if skip > len(data) {
		skip = len(data)
	}
	end := skip + pagesPerList
	if end > len(data) {
		end = len(data)
	}
	return data[skip:end]
}

func LeadersHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}
	usersTotal := len(models.Users)

	PageNum, err := strconv.Atoi(page)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	if PageNum > usersTotal/pagesPerList || PageNum < 1 {
		helpers.Return400(&w, helpers.ErrorSet{`Invalid "Page" Value`})
		return
	}

	userSlice := make([]interface{}, 0, usersTotal)
	for _, v := range models.Users {
		userSlice = append(userSlice, v)
	}

	sort.Slice(userSlice, func(i, j int) bool { return userSlice[i].(models.User).Won > userSlice[j].(models.User).Won })
	paginatedSlice := Paginate(userSlice, int(pagesPerList*(PageNum-1)))
	var pageSlice []models.User
	for _, user := range paginatedSlice {
		pageSlice = append(pageSlice, user.(models.User))
	}

	ResponseData, _ := json.Marshal(LeaderBoard{int(usersTotal / 4), int(PageNum), pageSlice})
	_, err = w.Write(ResponseData)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}
