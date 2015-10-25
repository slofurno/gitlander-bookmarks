package main

import (
	"testing"
)

func TestDownSample(t *testing.T) {

	conn := newUserConnection()

	user1 := newUser()
	user2 := newUser()

	bookmark := &Bookmark{
		Id:  "hey",
		Url: "gitlander.com",
	}

	user1.AddConnection(conn)

	user2.AddSub(user1)
	user2.UpdateBookmark(bookmark)

	msg := <-conn.Inbox
	t.Log(string(msg))

}
