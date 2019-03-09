package router

import (
	"bytes"
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"github.com/gorilla/mux"
	"net/http"
)

func GetRouter()(router *mux.Router){
	router = mux.NewRouter()
	router.HandleFunc("/register", GetRegisterHandler).Methods("GET") //TODO DELETE
	router.HandleFunc("/auth", GetAuthHandler).Methods("GET") //TODO DELETE
	router.HandleFunc("/register", RegisterHandler).Methods("POST")
	router.HandleFunc("/auth", AuthHandler).Methods("POST")
	router.HandleFunc("/profile", ProfileHandler).Methods("GET")
	//router.HandleFunc("/profile", ProfileUpdateHandler).Methods("PATCH")

	return
}

func GetRegisterHandler(w http.ResponseWriter, r *http.Request){
	_, err := w.Write([]byte(
		`
			<!DOCTYPE html>
			<html lang="ru">
			<head>
    			<meta charset="UTF-8">
    			<title>Регистрация</title>
			</head>
			<body>
				<h1>Регистрация</h1>
				<form name="register-form" method="post">
    				<label for="email">Email </label>
    				<input type="text" name="email" id="email" required>
    				<br>
    				<label for="password">Password </label>
    				<input type="password" name="password" id="password" required>
					<br>
    				<label for="nickname">Nickname </label>
    				<input type="text" name="nickname" id="nickname" required>
					<br>
					<button type="submit">Зарегистрироваться</button>
				</form>
			</body>
			</html>
        `))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetAuthHandler(w http.ResponseWriter, r *http.Request){
	_, err := w.Write([]byte(
		`
			<!DOCTYPE html>
			<html lang="ru">
			<head>
    			<meta charset="UTF-8">
    			<title>Авторизация</title>
			</head>
			<body>
				<h1>Авторизация</h1>
				<form name="auth-form">
    				<label for="email">Email </label>
    				<input type="text" name="email" id="email" required>
    				<br>
    				<label for="password">Password </label>
    				<input type="password" name="password" id="password" required>
					<br>
					<button type="submit" formmethod="post">Войти</button>
				</form>
			</body>
			</html>
        `))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	str := r.Form.Get("email")
	fmt.Println(str)
	_, userExist := models.Users[r.Form.Get("email")]
	if userExist {
		_, err := w.Write([]byte(
			`
			<!DOCTYPE html>
			<html lang="ru">
			<head>
    			<meta charset="UTF-8">
    			<title>Неудача</title>
			</head>
			<body>
				<h1>Ошибка: пользователь с таким e-mail уже зарегистрирован</h1>
			</body>
			</html>
        `))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	salt, err := MakeSalt()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user := models.User{
		ID:        models.MakeID(),
		Email:     r.Form.Get("email"),
		Salt:      salt,
		ProfileID: models.MakeID(),
		BestScore: 0,
	}
	user.Password = MakePasswordHash(r.Form.Get("password"), user.Salt)

	profile := models.Profile{
		ID:       user.ProfileID,
		Nickname: r.Form.Get("nickname"),
		AvatarID: models.MakeID(),
	}

	avatar := models.Avatar{
		ID: profile.AvatarID,
		Path: "img/default_avatar.jpg",
	}

	models.Profiles[profile.ID] = profile
	models.Avatars[avatar.ID] = avatar
	defer func() {
		//Пользователь уже успешно создан, поэтому его в любом случае следует добавить в БД
		//Однако, с ним ещё можно произвести полезную работу, которая может вызвать ошибки
		models.Users[user.Email] = user
		models.UserKeyPairs[user.ID] = user.Email
	}()

	sessionCookie, err := MakeSession(&user)//Заводим для пользователя новую сессию
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &sessionCookie)

	_, err = w.Write([]byte(
		`
			<!DOCTYPE html>
			<html lang="ru">
			<head>
    			<meta charset="UTF-8">
    			<title>Успех</title>
			</head>
			<body>
				<h1>Пользователь успешно зарегистрирован!</h1>
			</body>
			</html>
        `))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func AuthHandler(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, found := models.Users[r.Form.Get("email")]
	password := MakePasswordHash(r.Form.Get("password"), user.Salt)
	if !found || bytes.Compare(password, user.Password) != 0 {
		_, err := w.Write([]byte(
			`
			<!DOCTYPE html>
			<html lang="ru">
			<head>
    			<meta charset="UTF-8">
    			<title>Неудача</title>
			</head>
			<body>
				<h1>Ошибка: неправильно введён логин или пароль</h1>
			</body>
			</html>
        `))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	sessionCookie, err := MakeSession(&user)//Заводим для пользователя новую сессию
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &sessionCookie)
	_, err = w.Write([]byte(
		`
			<!DOCTYPE html>
			<html lang="ru">
			<head>
    			<meta charset="UTF-8">
    			<title>Успех</title>
			</head>
			<body>
				<h1>Пользователь успешно авторизован!</h1>
			</body>
			</html>
        `))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusUnauthorized)
		fmt.Println("Error while loading cookie")
		return
	}
	user, err := Authorize(sessionCookie.Value)
	if err != nil {
		http.Redirect(w, r, "/auth", http.StatusUnauthorized)
		fmt.Println("Error while doing auth")
		return
	}
	profile, found := models.Profiles[user.ProfileID]
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	avatar, found := models.Avatars[profile.AvatarID]
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write([]byte(
		`
			<!DOCTYPE html>
			<html lang="ru">
			<head>
    			<meta charset="UTF-8">
    			<title>Профиль</title>
			</head>
			<body>
				<h1>Hello, `+profile.Nickname+`</h1>
				<h2>Your email is: `+user.Email+`</h2>
				<img src="`+avatar.Path+`" alt="Avatar">
			</body>
			</html>
        `))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
