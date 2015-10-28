package main

import (
	"github.com/nu7hatch/gouuid"
)

type BookmarkRequest struct {
	User string `json:"user"`
	Url  string
	Tags []string
}

type Bookmark struct {
	Id   string
	Url  string
	Tags map[string]bool
}

type BookmarkResponse struct {
	User      string
	Bookmarks []*Bookmark
}

type BookmarkDump struct {
	User      string
	Bookmarks map[string]*Bookmark
}

func newBookmark(url string) *Bookmark {
	u, _ := uuid.NewV4()
	return &Bookmark{
		Id:   u.String(),
		Url:  url,
		Tags: map[string]bool{},
	}
}

type subChange struct {
	User *User
	Type string
}

type UserConnection struct {
	Socket *WebSocket
}

func newUserConnection() *UserConnection {
	return &UserConnection{
	//Inbox: make(chan []byte, 10),
	}
}
