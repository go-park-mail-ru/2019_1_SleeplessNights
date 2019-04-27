package message

type Message struct {
	Type    string      `json:"post"`
	Payload interface{} `json:"payload"`
}

// Messages got from frontend
type PostMessage struct {
	Text       string `json:"text"`
	AvatarPath string `json:"avatar_path"`
}

type ScrollMessage struct {
	Since int `json:"since"`
}

type LeaveMessage struct {
}
