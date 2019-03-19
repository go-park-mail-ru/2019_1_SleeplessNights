package handlers

import "net/http"

func ForumCreate(w http.ResponseWriter, r *http.Request) {
	//Создание форума
}

func ForumSlugCreate(w http.ResponseWriter, r *http.Request) {
	//Создание ветки
}

func ForumSlugDetails(w http.ResponseWriter, r *http.Request) {
	//Получение информации о форуме
}

func ForumSlugThreads(w http.ResponseWriter, r *http.Request) {
	//Список ветвей обсуждения форума
}

func ForumSlugUsers(w http.ResponseWriter, r *http.Request) {
	//Пльзователи данного форума
}
