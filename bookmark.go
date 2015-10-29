package main

import (
	"encoding/json"
	"fmt"
)

type BookmarkRequest struct {
	User string `json:"user"`
	Url  string
	Tags []string
}

type User struct {
	Id string
}

type UserConnection struct {
	handles       map[string]func()
	subscriptions *Collection
	onstop        func()
	outbox        chan []byte
	socket        *WebSocket
}

type Bookmark struct {
	Id    string
	Owner string
	Url   string
}

func newUserConnection(userSubs *Collection, socket *WebSocket) *UserConnection {
	//user := &UserConnection{bookmarks: newCollection(), subscriptions: newCollection()}

	self := &UserConnection{subscriptions: userSubs, handles: map[string]func(){}, socket: socket}

	subadded := func(key string, value interface{}) {

		var ok bool
		var userinfo *userInfo

		if userinfo, ok = userInfos[key]; !ok {
			fmt.Println("invalid collection key", key)
			return
		}

		added := func(key string, value interface{}) {

			j, err := json.Marshal(value)

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			fmt.Println("sending json:", j)
			socket.Write(j)
		}

		changed := func(key string, value interface{}) {

		}

		removed := func(key string, value interface{}) {
			fmt.Println("removed bookmark", value)
		}

		callback := &Callback{added, changed, removed}

		self.handles[key] = userinfo.bookmarks.ObserveChanges(callback)

		j, err := json.Marshal(value)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		fmt.Println("added: ", string(j))
	}

	subchanged := func(key string, value interface{}) {

		fmt.Println("tevs")
	}

	subremoved := func(key string, value interface{}) {

		onstop := self.handles[key]
		onstop()
		//check if we are subbed to key, if so, dlete from our map + call onstop()
		fmt.Println("added: ", key)
	}

	subcallback := &Callback{subadded, subchanged, subremoved}

	self.onstop = userSubs.ObserveChanges(subcallback)

	return self
}
