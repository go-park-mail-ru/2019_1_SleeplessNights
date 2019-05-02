package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

const limit = 10

func LeadersHandler(w http.ResponseWriter, r *http.Request) {
	sinceStr := r.URL.Query().Get("since")
	if sinceStr == "" {
		sinceStr = "0"
	}

	since, err := strconv.ParseUint(sinceStr, 10, 32)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	leaders, err := userManager.GetLeaderBoardPage(context.Background(),
		&services.PageData{
			Since: since,
			Limit: limit,
		})
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	data, err := json.Marshal(leaders)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}
}
