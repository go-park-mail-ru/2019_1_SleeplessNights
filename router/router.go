package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type User struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	SessionID uint   `json:"session-id"`
	ProfileID uint   `json:"profile-id"`
}

type Session struct {
	ID    uint
	Token string
	ExpiresAt time.Time
}

type Profile struct {
	ID       uint   `json:"id"`
	Nickname string `json:"nickname"`
	AvatarID uint   `json:"avatar-id"`
}

type Avatar struct {
	ID   uint
	Path string
}

var users []User
var sessions []Session
var profiles []Profile
var avatars []Avatar

func init() {
	//TODO init slices
}

func GetRouter()(router *mux.Router){
	router = mux.NewRouter()
	//router.HandleFunc("/register", registerHandler).Methods("PUSH")
	router.HandleFunc("/profile", ProfileCreationHandler).Methods("POST")
	router.HandleFunc("/profile", ProfileUpdateHandler).Methods("PATCH")
	router.HandleFunc("/profile", ProfileHandler).Methods("GET")
	router.HandleFunc("/profile/{username}", UserProfileHandler).Methods("GET")

	return
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("You are so pretty today!"))
	if err != nil {
		w.WriteHeader(500)
		return
	}
}

func ProfileCreationHandler(w http.ResponseWriter, r *http.Request)  {}
func ProfileUpdateHandler(w http.ResponseWriter, r *http.Request)  {}
func UserProfileHandler(w http.ResponseWriter, r *http.Request)  {}

/*func registerHandler(w http.ResponseWriter, r *http.Request){

}*/