package handlers

import (
	"github.com/gorilla/websocket"
	"net/http"
)

func UpgradeWs(w http.ResponseWriter, r *http.Request) {

	//Read cookie, Authorize?
	upgrader := websocket.Upgrader{}
	_, err := upgrader.Upgrade(w, r, http.Header{"Upgrade": []string{"websocket"}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//create player?

}
