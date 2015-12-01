package main

import (
	"encoding/json"
	"fmt"
)

type DataUnion struct {
	UserId   string
	Bookmark *Bookmark
	Sub      string
	Token    string
	Name     string
}

type DataStore struct {
}

//TODO move dup code from handler here
func (s *DataStore) UpsertBookmark(userinfo *userInfo, bookmark *Bookmark) {

	//TODO:if we update bookmarks here were gonna dup tags
	for _, tag := range bookmark.Tags {
		userinfo.summary[tag] += 1
	}

	userinfo.bookmarks.Add(bookmark.Id, bookmark)

	data := &DataUnion{
		UserId:   userinfo.userid,
		Bookmark: bookmark,
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	if db.init {
		db.bookmarks.WriteRecord(b)
	}

}

func (s *DataStore) AddUser(userinfo *userInfo, token string) {

	if _, ok := userInfos[userinfo.userid]; !ok {
		userinfo.subscriptions.Add(userinfo.userid, userinfo.name)
		userInfos[userinfo.userid] = userinfo
	}

	userTokens[token] = userinfo.userid

	data := &DataUnion{
		UserId: userinfo.userid,
		Token:  token,
		Name:   userinfo.name,
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	if db.init {
		db.users.WriteRecord(b)
	}

}

//TODO: prolly need a real type instead of storing sub name
func (s *DataStore) AddSubscription(userinfo *userInfo, sub string, name string) {

	userinfo.subscriptions.Add(sub, name)

	data := &DataUnion{
		UserId: userinfo.userid,
		Sub:    sub,
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	if db.init {
		db.subscriptions.WriteRecord(b)
	}
}
