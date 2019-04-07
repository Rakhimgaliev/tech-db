package models

import "time"

type Post struct {
	Id       int64     `json:"id,omitempty"`
	Parent   int64     `json:"parent,omitempty"`
	Author   string    `json:"author"`
	Message  string    `json:"message"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	Thread   int32     `json:"thread,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	// Created string  `json:"created,omitempty"`
	Path []int64 `json:"-"`
}

// type PostFull struct {
// 	Post   *Post   `json:"post,omitempty"`
// 	Author *User   `json:"author,omitempty"`
// 	Thread *Thread `json:"thread,omitempty"`
// 	Forum  *Forum  `json:"forum,omitempty"`
// }

// type PostUpdate struct {
// 	Message string `json:"message,omitempty"`
// }

type Posts []*Post
