package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/faker"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"math"
	"net/http"
	"net/http/httptest"
	"sort"
	"strconv"
	"testing"
)

func TestLeadersHandlerSuccessful(t *testing.T) {

	faker.CreateFakeData(UserCounter)
	usersTotal := len(models.Users)

	userSlice := make([]interface{}, 0, usersTotal)
	for _, v := range models.Users {
		userSlice = append(userSlice, v)
	}
	sort.Slice(userSlice, func(i, j int) bool { return userSlice[i].(models.User).Won > userSlice[j].(models.User).Won })

	pagesTotal := int(math.Ceil(float64(UserCounter / PagesPerList)))

	for i := 1; i <= pagesTotal; i++ {

		paginatedSlice := Paginate(userSlice, int(PagesPerList*(i-1)))
		var pageSlice []models.User
		for _, user := range paginatedSlice {
			pageSlice = append(pageSlice, user.(models.User))
		}

		expected := fmt.Sprintf("{\"pages_total\":%d,\"page\":%d,\"data\":[", pagesTotal, i)
		for ii, user := range pageSlice {
			str := fmt.Sprintf("{\"email\":\"%s\",\"won\":%d,\"lost\":%d,\"play_time\":%d,\"nickname\":\"%s\",\"avatar_path\":\"%s\"}",
				user.Email, user.Won, user.Lost, user.PlayTime, user.Nickname, user.AvatarPath)
			if ii == len(pageSlice)-1 {
				expected = expected + str
			} else {
				expected = expected + str + ","
			}
		}
		expected = expected + "]}"

		req := httptest.NewRequest(http.MethodGet, ApiLeader, nil)
		qq := req.URL.Query()
		qq.Add("page", strconv.Itoa(i))
		req.URL.RawQuery = qq.Encode()

		resp := httptest.NewRecorder()

		http.HandlerFunc(LeadersHandler).ServeHTTP(resp, req)
		if status := resp.Code; status == http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: %v\nhandler can't write into responce \n",
				status)
		} else {
			if status := resp.Code; status != http.StatusOK {
				t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
					status, http.StatusOK)
			}

			if resp.Body.String() != expected {
				t.Errorf("\nhandler returned unexpected body:\ngot %v\nwant %v\n",
					resp.Body.String(), expected)
			}
		}
	}
}

func TestLeadersHandlerUnsuccessfulWithWrongValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, ApiLeader, nil)
	qq := req.URL.Query()
	qq.Add("page", "aa")
	req.URL.RawQuery = qq.Encode()

	resp := httptest.NewRecorder()

	http.HandlerFunc(LeadersHandler).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: %v\nhandler can't write into responce \n",
			status)
	} else {
		if status := resp.Code; status != http.StatusNotFound {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n",
				status, http.StatusNotFound)
		}
	}
}

func TestLeadersHandlerUnsuccessfulWithWrongPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, ApiLeader, nil)
	qq := req.URL.Query()
	qq.Add("page", strconv.Itoa(1000))
	req.URL.RawQuery = qq.Encode()

	resp := httptest.NewRecorder()

	http.HandlerFunc(LeadersHandler).ServeHTTP(resp, req)

	if status := resp.Code; status == http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: %v\nhandler can't write into responce \n",
			status)
	} else {
		if status := resp.Code; status != http.StatusBadRequest {
			t.Errorf("\nhandler returned wrong status code:\ngot %v\nwant %v\n%s",
				status, http.StatusBadRequest, resp.Body.String())
		}

		expected := `{"email":"","password":"","password2":"","nickname":"","error":["Invalid \"Page\" Value"]}`
		if resp.Body.String() != expected {
			t.Errorf("\nhandler returned unexpected body:\ngot %v\nwant %v\n",
				resp.Body.String(), expected)
		}
	}
}
