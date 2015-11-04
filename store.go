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

func (s *DataStore) AddBookmark(userinfo *userInfo, bookmark *Bookmark) {

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
	database.WriteRecord(b)

}

func (s *DataStore) AddUser(userinfo *userInfo) {

	userinfo.subscriptions.Add(userinfo.userid, userinfo.name)
	userTokens[userinfo.token] = userinfo.userid
	userInfos[userinfo.userid] = userinfo

	data := &DataUnion{
		UserId: userinfo.userid,
		Token:  userinfo.token,
		Name:   userinfo.name,
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	database.WriteRecord(b)

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

	database.WriteRecord(b)
}
