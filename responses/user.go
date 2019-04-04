package responses

import "time"

type User struct {
	IsNew   bool      `json:"-"`
	ID      uint64    `json:"id"`
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Forum   string    `json:"forum"`
	Message string    `json:"message"`
	Votes   uint32    `json:"votes"`
	Slug    string    `json:"slug"`
	Created time.Time `json:"created"`
}
