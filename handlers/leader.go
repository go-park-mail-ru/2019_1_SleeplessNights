package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/database"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"math"
	"net/http"
	"strconv"
)

const (
	PagesPerList = 10
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
	end := skip + PagesPerList
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
	usersTotal, err := database.GetInstance().GetLenUsers()
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	if usersTotal == 0 {
		_, err := w.Write([]byte(`{"pages_total":0,"page":1,"data":[]}`))
		if err != nil {
			helpers.Return500(&w, err)
			return
		}
	}

	PageNum, err := strconv.Atoi(page)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	pagesTotal := int(math.Ceil(float64(usersTotal) / PagesPerList))
	if PageNum > pagesTotal || PageNum < 1 {
		helpers.Return400(&w, helpers.ErrorSet{`Invalid "Page" Value`})
		return
	}

	userSlice := make([]interface{}, 0, usersTotal)
	users, err := database.GetInstance().GetUsers()
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	for _, v := range users {
		userSlice = append(userSlice, v)
	}

	paginatedSlice := Paginate(userSlice, int(PagesPerList*(PageNum-1)))
	var pageSlice []models.User
	for _, user := range paginatedSlice {
		pageSlice = append(pageSlice, user.(models.User))
	}

	ResponseData, _ := json.Marshal(LeaderBoard{int(math.Ceil(float64(usersTotal / PagesPerList))), int(PageNum), pageSlice})
	_, err = w.Write(ResponseData)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}
