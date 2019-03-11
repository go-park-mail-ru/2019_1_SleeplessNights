package router

import (
	"encoding/json"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	DomainsCORS     = "https://sleeples-nights--frontend.herokuapp.com"
	MethodsCORS     = "GET POST PATCH OPTIONS"
	CredentialsCORS = "true"
	HeadersCORS     = "X-Requested-With, Content-type, User-Agent, Cache-Control, Cookie, Origin, Accept-Encoding, Connection, Host, Upgrade-Insecure-Requests, User-Agent, Referer, Access-Control-Request-Method, Access-Control-Request-Headers"
)

const (
	AvatarPrefix = "static/img/"
	MaxPhotoSize = 2*1024*1024
)

func GetRouter()(router *mux.Router){
	router = mux.NewRouter()
	router.HandleFunc("/api/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/api/auth", AuthHandler).Methods("POST")//.Headers("Referer")
	router.HandleFunc("/api/profile", func(w http.ResponseWriter, r *http.Request) {SetBasicHeaders(&w); w.WriteHeader(http.StatusOK)}).Methods("OPTIONS")
	router.HandleFunc("/api/profile", ProfileHandler).Methods("GET")
	//router.HandleFunc("/profile", ProfileUpdateHandler).Methods("PATCH", "OPTIONS")
	//router.HandleFunc("/api/leaders", LeadersHandler).Methods("GET")
	router.HandleFunc("/img/{path}", func(w http.ResponseWriter, r *http.Request) {SetBasicHeaders(&w); w.WriteHeader(http.StatusOK)}).Methods("OPTIONS")
	router.HandleFunc("/img/{path}", ImgHandler).Methods("GET") //.Queries("path")
	return
}


func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	SetBasicHeaders(&w)
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
	_, err = w.Write([]byte("{}"))
	if err != nil {
		Return500(&w, err)
		return
	}
}


func AuthHandler(w http.ResponseWriter, r *http.Request){
	SetBasicHeaders(&w)
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
	_, err = w.Write([]byte("{}"))
	if err != nil {
		Return500(&w, err)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	SetBasicHeaders(&w)
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := Authorize(sessionCookie.Value)
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
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

/*func ProfileUpdateHandler(w http.ResponseWriter, r *http.Request) {
	SetBasicHeaders(&w)
	err := r.ParseMultipartForm(MaxPhotoSize)
	if err != nil {
		formErrorMessages := ErrorSet{
			FormParsingErrorMsg,
			err.Error(),
		}
		Return400(&w, formErrorMessages)
		return
	}

	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := Authorize(sessionCookie.Value)
	if err != nil {
		r.Header.Add("Referer", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	requestErrors, isValid, err := ValidateUpdateProfileRequest(r, user)
	if err != nil {
		Return500(&w, err)
	}
	if !isValid {
		Return400(&w, requestErrors)
		return
	}

	user.Nickname = r.MultipartForm.Value["nickname"][0]
	models.Users[user.Email] = user
	newAvatar, err := r.MultipartForm.File["avatar"][0].Open()
	if err != nil {
		Return500(&w, err)
		return
	}
	file, err :=
}*/

func ImgHandler(w http.ResponseWriter, r *http.Request) {
	SetBasicHeaders(&w)
	vars := mux.Vars(r)
	pathToFile, found := vars["path"]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	path := AvatarPrefix + pathToFile//r.URL.Query().Get("path")
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

func LeadersHandler(w http.ResponseWriter, r *http.Request)  {

}

func SetBasicHeaders(w *http.ResponseWriter){
	(*w).Header().Set("Access-Control-Allow-Origin", DomainsCORS)
	(*w).Header().Set("Access-Control-Allow-Credentials", CredentialsCORS)
	(*w).Header().Set("Access-Control-Allow-Methods", MethodsCORS)
	(*w).Header().Set("Access-Control-Allow-Headers", HeadersCORS)
	(*w).Header().Set("Content-type", "application/json")
}