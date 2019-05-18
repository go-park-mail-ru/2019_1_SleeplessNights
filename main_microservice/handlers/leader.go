package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

var limit = uint64(config.GetInt("main_ms.pkg,handlers.leaderboard_page_len"))

func LeadersHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "0"
	}

	since, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		helpers.Return500(&w, err)
		return
	}

	leaders, err := userManager.GetLeaderBoardPage(context.Background(),
		&services.PageData{
			Since: since * limit,
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
