package handlers

import "net/http"

func ThreadCreate (w http.ResponseWriter, r *http.Request) {
	//Создание новых постов
}

func ThreadGetDetails (w http.ResponseWriter, r *http.Request) {
	//Получение информации о ветке обсуждения
}

func ThreadSetDetails (w http.ResponseWriter, r *http.Request) {
	//Обновление ветки
}

func ThreadPosts (w http.ResponseWriter, r *http.Request) {
	//Сообщение данной ветки обсуждения
}

func ThreadVote (w http.ResponseWriter, r *http.Request) {
	//Проголосовать за ветвь обсуждения
}