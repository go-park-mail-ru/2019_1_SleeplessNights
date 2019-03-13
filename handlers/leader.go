package handlers

import (
	"encoding/json"
	_ "encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/gorilla/mux"
	"net/http"
	"sort"
	"strconv"
)

const (
	pagesPerList = 4
)

type LeaderBoard struct {
	Pages_total  int           `json:"pages_total",string`
	CurrenPage   int           `json:"page"",string`
	LeadersSlice []models.User `json:"data"`
}

func Paginate(x []models.User, skip int, size int) []models.User {
	if skip > len(x) {
		skip = len(x)
	}
	end := skip + pagesPerList
	if end > len(x) {
		end = len(x)
	}
	return x[skip:end]
}

func LeadersHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	page, found := vars["page"]

	if !found {
		page = "1"
	}

	pagenum, err := strconv.ParseInt(page, 10, 32)

	if err != nil {
		helpers.Return400(&w, helpers.ErrorSet{`Invalid "Page" Value`})
		return
	}
	usersTotal := len(models.Users)
	if found && (pagenum > int64(usersTotal/pagesPerList) || pagenum < 1) {
		helpers.Return400(&w, helpers.ErrorSet{`Invalid "Page" Value`})
		return
	}
	UserSlice := make([]models.User, usersTotal)
	i := 0
	for _, v := range models.Users {
		UserSlice[i] = v
		i += 1
	}
	sort.Slice(UserSlice, func(i, j int) bool { return UserSlice[i].Won > UserSlice[j].Won })
	PageSlice := Paginate(UserSlice, int(pagesPerList*(pagenum-1)), len(UserSlice))
	ResponseData, _ := json.Marshal(LeaderBoard{int(usersTotal / 4), int(pagenum), PageSlice})
	w.Header().Set("Content-type", "application/json")
	w.Write(ResponseData)

}
