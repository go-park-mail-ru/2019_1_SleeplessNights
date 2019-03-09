package router

import "net/http"

type User struct {
	ID        uint
	Email     string
	Password  []byte//Длина хеш-суммы по алгоритму SHA512
	Salt      []byte
	SessionID uint
	ProfileID uint
	BestScore uint
}

type UserPk struct {
	ID    uint
	Email string
}

/*type Session struct {
	ID        uint
	Token     []byte
	CreatedBy string
	CreatedAt time.Time
	ExpiresAt time.Time
}*/

type Profile struct {
	ID       uint
	Nickname string
	AvatarID uint
}

type Avatar struct {
	ID   uint
	Path string
}

var IDSource uint
var users map[string]User
var userKeyPairs map[uint]string
var profiles map[uint]Profile
var avatars map[uint]Avatar
var secret []byte

func init() {
	users = make(map[string]User, 0)
	userKeyPairs = make(map[uint]string, 0)
	profiles = make(map[uint]Profile, 0)
	avatars = make(map[uint]Avatar, 0)

	err := os.Setenv("SaltLen", "16")//16 байт (128 бит), как в современных UNIX системах
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("SessionLifeLen", time.Hour.String())
	if err != nil {
		log.Fatal(err)
	}
	err = os.Setenv("ServerID", "MyConfidantServer")
	if err != nil {
		log.Fatal(err)
	}
	secretFile, err := os.Open("secret")
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fscanln(secretFile, &secret)
	if err != nil {
		return
	}
}

func MakeID() uint {
	//TODO make thread-safe
	IDSource++
	return IDSource
}

func MakeSalt()(salt []byte, err error) {
	saltLen, err := strconv.Atoi(os.Getenv("SaltLen"))
	if err != nil {
		return
	}
	salt = make([]byte, saltLen)
	rand.Read(salt)//Заполняем слайс случайными значениями по всей его длине
	return
}

func MakePasswordHash(password string, salt []byte)(hash []byte){
	saltedPassword:= bytes.Join([][]byte {[]byte(password), salt}, nil)
	hashedPassword := sha512.Sum512(saltedPassword)//sha512 возвращает массив, а слайс можно взять только по addressable массиву
	hash = hashedPassword[0:]
	return
}

func MakeSession(user *User)(sessionCookie http.Cookie, err error){
	signer := jwt.NewHMAC(jwt.SHA512, secret)
	header := jwt.Header{}
	sessionLifeLen, err := time.ParseDuration(os.Getenv("SessionLifeLen"))
	if err != nil {
		return
	}
	expiresAt := time.Now().Add(sessionLifeLen)
	payload := jwt.Payload{
		ExpirationTime: expiresAt.Unix(),
		JWTID: strconv.FormatUint(uint64(user.ID), 10),
	}
	token, err := jwt.Sign(header, payload, signer)
	if err != nil {
		return
	}

	//Создаём cookie сессии
	sessionCookie = http.Cookie{
		Name:     "session_token",
		Value:    string(token),
		Expires:  expiresAt,
		HttpOnly: true,
	}
	return
}

func Authorize(sessionToken string)(*User, error){
	rawToken, err := jwt.Parse([]byte(sessionToken))
	if err != nil {
		fmt.Println("Error while parsing token")
		return nil, err
	}
	verifier := jwt.NewHMAC(jwt.SHA512, secret)
	err = rawToken.Verify(verifier)
	if err != nil {
		fmt.Println("Error while verifying token")
		return nil, err
	}
	payload := jwt.Payload{}
	_, err = rawToken.Decode(&payload)
	if err != nil {
		fmt.Println("Error while decoding token")
		fmt.Println(err)
		return nil, err
	}
	expValidator := jwt.ExpirationTimeValidator(time.Now(), true)
	err = payload.Validate(expValidator)
	if err != nil {
		fmt.Println("Error while validating token")
		return nil, err
	}
	userID, err := strconv.ParseUint(payload.JWTID, 10, 64)
	if err != nil {
		return nil, err
	}
	user, found := users[userKeyPairs[uint(userID)]]
	if !found {
		fmt.Println("Error: user not found in the database")
		return nil, errors.New("error: There are no token's owner in database")
	}
	return &user, nil
}

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
	_, userExist := users[r.Form.Get("email")]
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
	user := User{ID:        MakeID(),
		Email:     r.Form.Get("email"),
		Salt:      salt,
		ProfileID: MakeID(),
		BestScore: 0}
	user.Password = MakePasswordHash(r.Form.Get("password"), user.Salt)

	profile := Profile{ID:       user.ProfileID,
		Nickname: r.Form.Get("nickname"),
		AvatarID: MakeID()}

	avatar := Avatar{ID: profile.AvatarID,
		Path: "img/default_avatar.jpg"}

	profiles[profile.ID] = profile
	avatars[avatar.ID] = avatar
	defer func() {
		//Пользователь уже успешно создан, поэтому его в любом случае следует добавить в БД
		//Однако, с ним ещё можно произвести полезную работу, которая может вызвать ошибки
		users[user.Email] = user
		userKeyPairs[user.ID] = user.Email
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
	user, found := users[r.Form.Get("email")]
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
	profile, found := profiles[user.ProfileID]
	if !found {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	avatar, found := avatars[profile.AvatarID]
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
