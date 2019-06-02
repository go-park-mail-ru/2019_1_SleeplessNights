package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/main_microservice/handlers/helpers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/services"
	"golang.org/x/net/context"
	"net/http"
	"strconv"
)

func LeadersHandler(w http.ResponseWriter, r *http.Request) {
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "1"
	}

	since, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		logger.Errorf("Failed to parse page: %v", err.Error())
		helpers.Return400(&w,helpers.ErrorSet{err.Error()})
		return
	}

	leaders, err := userManager.GetLeaderBoardPage(context.Background(),
		&services.PageData{
			Since: since,
		})
	if err != nil {
		logger.Errorf("Failed to get leader board page: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}

	if len(leaders.Leaders) == 0 {
		logger.Errorf("Page didn't find")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	data, err := json.Marshal(leaders)
	if err != nil {
		logger.Errorf("Failed to marshal leaders: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		logger.Errorf("Failed to write response: %v", err.Error())
		helpers.Return500(&w, err)
		return
	}
}
