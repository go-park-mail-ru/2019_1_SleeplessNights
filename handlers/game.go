package handlers

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/TheGame/messge"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/logger"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

func UpgradeWs(w http.ResponseWriter, r *http.Request) {
	//Read cookie, Authorize?
	cookie, err := r.Cookie("auth")
	player_id, err := strconv.ParseUint(cookie.Value, 10, 32)
	if err != nil {
		logger.Error.Println("Ws Connection handler ")
	}
	if err != nil {
		fmt.Println(err)
	}

	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, http.Header{"Upgrade": []string{"websocket"}})
	fmt.Println("player %s connected to socket", cookie.Value)
	if err != nil {
		fmt.Println("Error during UpgradeWs", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	err = conn.WriteJSON(messge.Message{"INFO", "you've been connected to server"})
	if err != nil {
		fmt.Println(err)
	}

	game := TheGame.GetInstance()
	game.PlayByWebsocket(conn, player_id)

}
