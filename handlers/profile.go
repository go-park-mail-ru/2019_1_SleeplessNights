package handlers

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/router"
	"net/http"
)


func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	router.SetBasicHeaders(&w)
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte("{}"))
		if err != nil {
			router.Return500(&w, err)
			return
		}
		return
	}
	user, err := router.Authorize(sessionCookie.Value)
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		_, err = w.Write([]byte("{}"))
		if err != nil {
			router.Return500(&w, err)
			return
		}
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		router.Return500(&w, err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		router.Return500(&w, err)
		return
	}
}

func ProfileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	router.SetBasicHeaders(&w)
	err := r.ParseMultipartForm(router.MaxPhotoSize)
	if err != nil {
		formErrorMessages := router.ErrorSet{
			router.FormParsingErrorMsg,
			err.Error(),
		}
		router.Return400(&w, formErrorMessages)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := router.Authorize(sessionCookie.Value)
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestErrors, isValid, err := router.ValidateUpdateProfileRequest(r, user)//TODO WRITE VALIDATOR
	if err != nil {
		router.Return500(&w, err)
	}
	if !isValid {
		router.Return400(&w, requestErrors)
		return
	}

	user.Nickname = r.MultipartForm.Value["nickname"][0]
	models.Users[user.Email] = user
	//TODO UPLOAD AVATAR
	newAvatar, err := r.MultipartForm.File["avatar"][0].Open()
	if err != nil {
		router.Return500(&w, err)
		return
	}
	file, err :=
}
