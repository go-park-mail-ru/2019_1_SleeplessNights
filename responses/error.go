package responses

import "encoding/json"

type Error struct {
	Message string `json:"message"`
}

func (err *Error)MarshalToJSON()([]byte, error) {
	return json.Marshal(err)
}
