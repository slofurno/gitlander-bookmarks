package main

import (
	"encoding/json"
	"github.com/nu7hatch/gouuid"
)

type User struct {
	Id          string
	Bookmarks   map[string]*Bookmark
	Subscribers map[string]*User

	Connections []*UserConnection
	Inbox       chan []byte
	subs        chan *subChange
	updates     chan *Bookmark

	listeners chan *UserConnection
}

func newUser() *User {

	u, _ := uuid.NewV4()

	user := &User{
		Id:          u.String(),
		Bookmarks:   make(map[string]*Bookmark),
		Subscribers: make(map[string]*User),
		Connections: make([]*UserConnection, 0),
		Inbox:       make(chan []byte, 10),
		subs:        make(chan *subChange, 10),
		updates:     make(chan *Bookmark, 10),
		listeners:   make(chan *UserConnection, 4),
	}

	go user.Publish()
	return user
}

func (s *User) Publish() {

	for {
		select {

		case connection := <-s.listeners:
			s.Connections = append(s.Connections, connection)

		case message := <-s.Inbox:

			for i := len(s.Connections) - 1; i >= 0; i-- {
				select {
				case s.Connections[i].Inbox <- message:
					//add a ping to make sure dropped conns get removed
				default:
					s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
				}
			}

		case subscriber := <-s.subs:
			if subscriber.Type == "add" {
				s.Subscribers[subscriber.User.Id] = subscriber.User
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

	j, _ := json.Marshal(bookmark)

	for _, user := range s.Subscribers {
		user.Inbox <- j
	}
}

func (s *User) AddSub(user *User) {
	s.subs <- &subChange{
		User: user,
		Type: "add",
	}
}

func (s *User) RemoveSub(user *User) {
	s.subs <- &subChange{
		User: user,
		Type: "remove",
	}
}
