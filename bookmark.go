package main

import (
	"encoding/json"
	"fmt"
	"github.com/slofurno/ws"
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
	Id          string   `json:"id"`
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
	socket        *ws.WebSocket
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

type BookmarkOp string

const (
	Add    BookmarkOp = "add"
	Delete BookmarkOp = "delete"
)

type BookmarkEvent struct {
	Type string
	Op   string
	Data string
}

func newUserConnection(userSubs *Collection, socket *ws.WebSocket) *UserConnection {

	self := &UserConnection{subscriptions: userSubs, handles: map[string]func(){}, socket: socket}

	subadded := func(key string, value interface{}) {

		var ok bool
		var userinfo *userInfo

		if userinfo, ok = userInfos[key]; !ok {
			fmt.Println("invalid collection key", key)
			return
		}

		//TODO find another way to get user info
		newsub := &userinfoDto{Name: userinfo.name, Userid: key}
		fmt.Println("sub:", userinfo.name, key)
		jj, _ := json.Marshal(newsub)
		socket.Write(jj)

		event := &BookmarkEvent{Type: "sub", Op: "add", Data: key}

		tevs, err := json.Marshal(event)

		if err != nil {
			fmt.Println(err.Error())
		} else {
			socket.Write(tevs)
		}

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
			bookmark, ok := value.(*Bookmark)
			if !ok {
				return
			}

			j, err := json.Marshal(bookmark)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("updated bookmark:", string(j))
			socket.Write(j)
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
			//TODO: probably delete onstop from map / use closure
			onstop()
		}

		//oldsub := &userinfoDto{Userid: key}
		//zz, err := json.Marshal(oldsub)

		event := &BookmarkEvent{Type: "sub", Op: "delete", Data: key}

		j, err := json.Marshal(event)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		socket.Write(j)
	}

	subcallback := &Callback{subadded, subchanged, subremoved}
	self.onstop = userSubs.ObserveChanges(subcallback)

	return self
}
