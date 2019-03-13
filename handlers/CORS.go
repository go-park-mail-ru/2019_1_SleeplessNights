package handlers

import (
	"net/http"
)

func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
