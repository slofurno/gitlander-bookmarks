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

type GithubUserResponse struct {
	Login string `json:"login"`
	Url   string `json:"url"`
	Name  string `json:"name"`
	Id    uint64 `json:"id"`
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
	//stores func to stop observing... can i move to closure?
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

		tevs := &userinfoDto{Name: userinfo.name, Userid: key}

		jj, _ := json.Marshal(tevs)
		socket.Write(jj)

		added := func(key string, value interface{}) {

			bookmark, ok := value.(*Bookmark)

			if !ok {
				fmt.Println("how is this not a bookmark?")
				return
			}

			j, err := json.Marshal(bookmark)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			socket.Write(j)
		}

		//TODO: consider how changing a bookmark's tags will affect our usersummary...
		changed := func(key string, value interface{}, old interface{}) {
			//
		}

		removed := func(key string, value interface{}) {
			fmt.Println("removed bookmark", value)
		}

		callback := &Callback{added, changed, removed}
		self.handles[key] = userinfo.bookmarks.ObserveChanges(callback)

	}

	subchanged := func(key string, value interface{}, old interface{}) {
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
