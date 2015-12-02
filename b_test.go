package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDownSample(t *testing.T) {
	/*
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
	*/
}

func TestRemove(t *testing.T) {
	j := []byte("{\"Id\":\"43534534\"}")
	bookmark := &Bookmark{}
	err := json.Unmarshal(j, bookmark)

	if err != nil {
		t.Error(err.Error())
	}

	if bookmark.Url == "" {
		fmt.Println("empty string")
	}
	fmt.Println(bookmark.Url)

}
