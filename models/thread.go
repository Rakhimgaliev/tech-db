package models

type Thread struct {
	Id      int    `json:"id,omitempty"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Forum   string `json:"forum,omitempty"`
	Message string `json:"message"`
	Votes   int    `json:"votes,omitempty"`
	Slug    string `json:"slug,omitempty"`
	// Created time.Time `json:"created,omitempty"`
	Created string `json:"created,omitempty"`
}

type ThreadUpdate struct {
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

type Threads []*Thread
