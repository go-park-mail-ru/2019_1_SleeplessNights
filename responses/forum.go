package responses

type Forum struct {
	PostsCount   int64 `json:"posts"`
	ForumSlug    string `json:"slug"`
	ThreadsCount int32 `json:"threads"`
	ForumTitle   string `json:"title"`
	UserNickname string `json:"user"`
	IsNew        bool   `json:"-"`
}

