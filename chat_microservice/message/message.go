package message

type Message struct {
	Title   string      `json:"title"`
	Payload interface{} `json:"payload"`
}

// Messages got from frontend
type PostMessage struct {
	Text       string `json:"text"`
	AvatarPath string `json:"avatarPath"`
}

type ScrollMessage struct {
	Since int `json:"since"`
}

type LeaveMessage struct {
}
