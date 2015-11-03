package main

import (
	"encoding/json"
	"fmt"
	"time"
)

func getCurrentTime() int64 {
	nanos := time.Now().UnixNano()
	return nanos / 1000000
}

type BookmarkRequest struct {
	//User        string   `json:"user"`
	Url         string   `json:"url"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
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
	Id          string
	Owner       string
	Url         string
	Description string
	Tags        []string
	Summary     string
	Time        int64
}

func newUserConnection(userSubs *Collection, socket *WebSocket) *UserConnection {

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
			//
		}

		removed := func(key string, value interface{}) {
			fmt.Println("removed bookmark", value)
		}

		callback := &Callback{added, changed, removed}
		self.handles[key] = userinfo.bookmarks.ObserveChanges(callback)

	}

	subchanged := func(key string, value interface{}) {
		fmt.Println("subchanged")
	}

	subremoved := func(key string, value interface{}) {
		if onstop, ok := self.handles[key]; ok {
			fmt.Println("unsubbed from: ", key)
			onstop()
		}
	}

	subcallback := &Callback{subadded, subchanged, subremoved}
	self.onstop = userSubs.ObserveChanges(subcallback)

	return self
}
