package responses

import (
	"encoding/json"
)

type Forum struct {
	PostsCount   uint64 `json:"posts"`
	ForumSlug    string `json:"slug"`
	ThreadsCount uint32 `json:"threads"`
	ForumTitle   string `json:"title"`
	UserNickname string `json:"user"`
}

func (forum *Forum)MarshalToJSON()([]byte, error) {
	return json.Marshal(forum)
}

