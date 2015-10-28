package main

import (
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
)

type User struct {
	Id            string
	Bookmarks     map[string]*Bookmark
	Subscribers   map[string]*User
	Subscriptions map[string]*User

	Connections []*UserConnection
	Inbox       chan []byte
	subs        chan *subChange
	updates     chan *Bookmark
	fullUpdates chan chan []byte

	listeners chan *UserConnection
	leavers   chan *UserConnection
}

func newUser() *User {

	u, _ := uuid.NewV4()

	user := &User{
		Id:            u.String(),
		Bookmarks:     make(map[string]*Bookmark),
		Subscribers:   make(map[string]*User),
		Subscriptions: make(map[string]*User),
		Connections:   []*UserConnection{},
		Inbox:         make(chan []byte, 128),
		subs:          make(chan *subChange, 16),
		updates:       make(chan *Bookmark, 32),
		listeners:     make(chan *UserConnection, 16),
		leavers:       make(chan *UserConnection, 16),
		fullUpdates:   make(chan chan []byte, 8),
	}

	go user.Publish()
	user.AddSub(user)

	return user
}

func (s *User) Publish() {

	for {
		select {

		case fullrequest := <-s.fullUpdates:

			dump := []*Bookmark{}

			for _, bookmark := range s.Bookmarks {
				dump = append(dump, bookmark)
			}

			response := &BookmarkResponse{
				User:      s.Id,
				Bookmarks: dump,
			}
			json, _ := json.Marshal(response)
			fullrequest <- json

		case connect := <-s.leavers:
			for i := len(s.Connections) - 1; i >= 0; i-- {
				if connect == s.Connections[i] {
					s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
				}
			}

		case connection := <-s.listeners:
			fmt.Println("adding listener")

			s.Connections = append(s.Connections, connection)
			world := make(chan (<-chan []byte), len(s.Subscriptions))

			for _, subscription := range s.Subscriptions {
				world <- subscription.GetBookmarks()
			}

			close(world)

			go func() {
				for result := range world {

					b := <-result
					fmt.Println(string(b))
					connection.Socket.Write(b)
				}

			}()

		case message := <-s.Inbox:
			fmt.Println("writing msg to connections: ", len(s.Connections))
			for i := len(s.Connections) - 1; i >= 0; i-- {
				s.Connections[i].Socket.Write(message)
				/*
					        select {
									case s.Connections[i].Inbox <- message:
										//add a ping to make sure dropped conns get removed
									default:
										s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
									}
				*/
			}

		case subscriber := <-s.subs:
			if subscriber.Type == "add" {
				fmt.Println("adding subscriber")
				s.Subscribers[subscriber.User.Id] = subscriber.User

				go func() {
					subscriber.User.Inbox <- <-s.GetBookmarks()
				}()

			} else if subscriber.Type == "sub" {
				s.Subscriptions[subscriber.User.Id] = subscriber.User
			} else if subscriber.Type == "unsub" {
				delete(s.Subscriptions, subscriber.User.Id)
			} else {
				delete(s.Subscribers, subscriber.User.Id)
			}

		case bookmark := <-s.updates:
			if bookmark.Url == "" {
				delete(s.Bookmarks, bookmark.Id)
			} else {
				s.Bookmarks[bookmark.Id] = bookmark
			}
			s.broadcast(bookmark)
		}
	}

}

func (s *User) AddConnection(connection *UserConnection) {
	s.listeners <- connection
}

func (s *User) UpdateBookmark(bookmark *Bookmark) {
	s.updates <- bookmark
}

func (s *User) broadcast(bookmark *Bookmark) {

	update := &BookmarkResponse{
		User:      s.Id,
		Bookmarks: []*Bookmark{bookmark},
	}

	j, _ := json.Marshal(update)

	for _, user := range s.Subscribers {
		user.Inbox <- j
	}
}

func (s *User) AddSub(user *User) {
	s.subs <- &subChange{
		User: user,
		Type: "add",
	}

	user.subs <- &subChange{
		User: s,
		Type: "sub",
	}
}

func (s *User) RemoveSub(user *User) {

	user.subs <- &subChange{
		User: s,
		Type: "unsub",
	}

	s.subs <- &subChange{
		User: user,
		Type: "remove",
	}
}

func (s *User) GetBookmarks() <-chan []byte {
	result := make(chan []byte, 1)
	s.fullUpdates <- result
	return result
}
