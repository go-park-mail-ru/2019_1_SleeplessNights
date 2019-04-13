package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
	"github.com/gorilla/websocket"
	"net/http"
)

func UpgradeWs(w http.ResponseWriter, r *http.Request) {

	//Read cookie, Authorize?
	cookie, err := r.Cookie("auth")

	if err != nil {
		fmt.Println(err)
	}

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, http.Header{"Upgrade": []string{"websocket"}})
	fmt.Println("player %s joined", cookie.Value)
	if err != nil {
		fmt.Println("Error during UpgradeWs", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	err = conn.WriteJSON(messge.Message{"INFO", nil})
	if err != nil {
		fmt.Println(err)
	}

	//create player?

}
