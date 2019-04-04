package responses

import (
	"fmt"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
}

func InternalError(msg string)(code int, errPtr *Error) {
	fmt.Println(msg)
	var err Error
	err.Message = msg
	errPtr = &err
	code = http.StatusInternalServerError
	return
}
