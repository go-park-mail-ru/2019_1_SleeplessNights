package router

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	DomainsCORS     = "*"
	MethodsCORS     = "GET POST PATCH OPTIONS"
	CredentialsCORS = "true"
	HeadersCORS     = "X-Requested-With, Content-type, User-Agent, Cache-Control"
)

const (
	AvatarPrefix = "../static/img/"
)

func GetRouter()(router *mux.Router){
	router = mux.NewRouter()
	router.HandleFunc("/api/register", RegisterHandler).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth", AuthHandler).Methods("POST", "OPTIONS")//.Headers("Referer")
	router.HandleFunc("/api/profile", ProfileHandler).Methods("GET", "OPTIONS")
	//router.HandleFunc("/profile", ProfileUpdateHandler).Methods("PATCH")
	router.HandleFunc("/api/avatar", AvatarHandler).Methods("GET, OPTIONS").Queries("path")
	//router.HandleFunc("/api/leaders", ProfileHandler).Methods("GET", "OPTIONS")
	return
}


func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	AllowCORS(&w)
	err := r.ParseForm()
	if err != nil {
		formErrorMessages := ErrorSet{
			FormParsingErrorMsg,
			err.Error(),
		}
		Return400(&w, formErrorMessages)
		return
	}

	requestErrors, isValid, err := ValidateRegisterRequest(r)
	if err != nil {
		Return500(&w, err)
	}
	if !isValid {
		Return400(&w, requestErrors)
		return
	}

	user := models.User{
		ID:        models.MakeID(),
		Email:     r.Form.Get("email"),
		Won:       0,
		Lost:      0,
		PlayTime:  0,
		Nickname: r.Form.Get("nickname"),
		AvatarPath: "default_avatar.jpg",
	}
	salt, err := MakeSalt()
	if err != nil {
		Return500(&w, err)
		return
	}
	user.Salt = salt
	user.Password = MakePasswordHash(r.Form.Get("password"), user.Salt)
	defer func() {
		//Пользователь уже успешно создан, поэтому его в любом случае следует добавить в БД
		//Однако, с ним ещё можно произвести полезную работу, которая может вызвать ошибки
		models.Users[user.Email] = user
		models.UserKeyPairs[user.ID] = user.Email//Пара ключей ID-email, чтобы юзера можно было найти 2-мя способами
	}()

	sessionCookie, err := MakeSession(user)
	if err != nil {
		Return500(&w, err)
		return
	}
	http.SetCookie(w, &sessionCookie)
	http.Redirect(w, r, "/profile", http.StatusFound)
}


func AuthHandler(w http.ResponseWriter, r *http.Request){
	AllowCORS(&w)
	err := r.ParseForm()
	if err != nil {
		formErrorMessages := ErrorSet{
			FormParsingErrorMsg,
			err.Error(),
		}
		Return400(&w, formErrorMessages)
		return
	}

	requestErrors, isValid, user, err := ValidateAuthRequest(r)
	if err != nil {
		Return500(&w, err)
	}
	if !isValid {
		Return400(&w, requestErrors)
		return
	}

	sessionCookie, err := MakeSession(user)
	if err != nil {
		Return500(&w, err)
		return
	}
	http.SetCookie(w, &sessionCookie)
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
}


func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	AllowCORS(&w)
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		http.Redirect(w, r, "/auth", http.StatusUnauthorized)
		return
	}
	user, err := Authorize(sessionCookie.Value)
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		http.Redirect(w, r, "/auth", http.StatusUnauthorized)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		Return500(&w, err)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		Return500(&w, err)
		return
	}
}

func AvatarHandler(w http.ResponseWriter, r *http.Request) {
	AllowCORS(&w)
	path := AvatarPrefix + strings.TrimPrefix(r.URL.Query().Get("path"), "/")
	_, err := os.Stat(path)
	if 	err !=nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	avatar, err := ioutil.ReadFile(path)
	if err != nil {
		Return500(&w, err)
		return
	}
	w.Header().Set("Content-type", http.DetectContentType(avatar))
	_, err = w.Write(avatar)
	if err != nil {
		Return500(&w, err)
		return
	}
}

func AllowCORS(w *http.ResponseWriter){
	(*w).Header().Set("Access-Control-Allow-Origin", DomainsCORS)
	(*w).Header().Set("Access-Control-Allow-Credentials", CredentialsCORS)
	(*w).Header().Set("Access-Control-Allow-Methods", MethodsCORS)
	(*w).Header().Set("Access-Control-Allow-Headers", HeadersCORS)
}