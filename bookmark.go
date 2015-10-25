package main

import (
	"github.com/nu7hatch/gouuid"
)

type Bookmark struct {
	Id   string
	Url  string
	Tags map[string]bool
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
