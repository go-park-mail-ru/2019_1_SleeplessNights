package handlers

import (
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/models"
	"net/http"
)

type Handler struct {
	handlerFunc func (w http.ResponseWriter, r *http.Request)
}

func (h *Handler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(w, r)
}

type AuthorizedHandler struct {
	User models.User
	handlerFunc func (user models.User, w http.ResponseWriter, r *http.Request)
}

func (h *AuthorizedHandler) ServeHTTP (w http.ResponseWriter, r *http.Request) {
	h.handlerFunc(h.User, w, r)
}