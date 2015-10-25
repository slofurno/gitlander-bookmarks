package main

import ()

type Bookmark struct {
	Id   string
	Url  string
	Tags map[string]bool
}

type subChange struct {
	User *User
	Type string
}

type UserConnection struct {
	Inbox chan []byte
}

func newUserConnection() *UserConnection {
	return &UserConnection{
		Inbox: make(chan []byte, 10),
	}
}
