package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"net/http"
	"sort"
	"strconv"
)

const (
	pagesPerList = 4
)

type LeaderBoard struct {
	Pages_total int           `json:"pages_total"`
	CurrenPage  int           `json:"page"`
	Slice       []models.User `json:"data"`
}

func Paginate(data []models.User, skip int, size int) []models.User {

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

	PageNum, err := strconv.ParseInt(page, 10, 32)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	if PageNum > int64(usersTotal/pagesPerList) || PageNum < 1 {
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
	PageSlice := Paginate(UserSlice, int(pagesPerList*(PageNum-1)), len(UserSlice))
	ResponseData, _ := json.Marshal(LeaderBoard{int(usersTotal / 4), int(PageNum), PageSlice})
	w.Header().Set("Content-type", "application/json")
	w.Write(ResponseData)

}
