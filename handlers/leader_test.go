package handlers_test

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/handlers"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"testing"
)

func TestLeadersHandler(t *testing.T) {

	path := "/api/leaders"

	faker.CreateFakeData(handlers.UserCounter)
	usersTotal := len(models.Users)

	userSlice := make([]interface{}, 0, usersTotal)
	for _, v := range models.Users {
		userSlice = append(userSlice, v)
	}
	sort.Slice(userSlice, func(i, j int) bool { return userSlice[i].(models.User).Won > userSlice[j].(models.User).Won })

	var pagesTotal int
	if handlers.UserCounter % handlers.PagesPerList == 0{
		pagesTotal = handlers.UserCounter/handlers.PagesPerList
	} else{
		pagesTotal = handlers.UserCounter/handlers.PagesPerList+1
	}

	for i := 1; i <= pagesTotal; i++ {

		paginatedSlice := handlers.Paginate(userSlice, int(handlers.PagesPerList*(i-1)))
		var pageSlice []models.User
		for _, user := range paginatedSlice {
			pageSlice = append(pageSlice, user.(models.User))
		}

		expected := fmt.Sprintf("{\"pages_total\":%d,\"page\":%d,\"data\":[",pagesTotal, i)
		for ii, user := range pageSlice{
			str := fmt.Sprintf("{\"email\":\"%s\",\"won\":%d,\"lost\":%d,\"play_time\":%d,\"nickname\":\"%s\",\"avatar_path\":\"%s\"}",
				user.Email, user.Won, user.Lost, user.PlayTime, user.Nickname, user.AvatarPath)
			if ii == len(pageSlice) - 1{
				expected = expected + str
			} else {
				expected = expected + str + ","
			}
		}
		expected = expected + "]}"

		req := httptest.NewRequest(http.MethodGet, path, nil)
		qq := req.URL.Query()
		qq.Add("page", strconv.Itoa(i))
		req.URL.RawQuery = qq.Encode()

		resp := httptest.NewRecorder()

		http.HandlerFunc(handlers.LeadersHandler).ServeHTTP(resp, req)

		if status := resp.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code:\n got %v\n want %v\n",
				status, http.StatusOK)
		}

		if resp.Body.String() != expected {
			t.Errorf("handler returned unexpected body:\n got %v\n want %v\n",
				resp.Body.String(), expected)
		}
	}
}
