package models

type User struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname"`
}

type UserUpdate struct {
	About    string `json:"about,omitempty"`
	Email    string `json:"email,omitempty"`
	Fullname string `json:"fullname,omitempty"`
}

type Users []*User
