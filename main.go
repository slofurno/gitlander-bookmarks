package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"net/http"
)

var userTokens = map[string]string{}
var userInfos = map[string]*userInfo{}

var secretKey []byte
var dataStore = &DataStore{}

var client_id string
var client_secret string

var globalSummary = map[string]int{}

var db struct {
	bookmarks     *Filebase
	subscriptions *Filebase
	users         *Filebase
	init          bool
}

type clientSecrets struct {
	Client_id     string
	Client_secret string
}

type userInfo struct {
	subscriptions *Collection
	bookmarks     *Collection
	summary       map[string]int
	userid        string
	name          string
}

type userinfoDto struct {
	Summary map[string]int
	Name    string
	Userid  string
}

type RequestContext struct {
	isAuthed bool
	userinfo *userInfo
	user     string
}

func newUserInfo() *userInfo {

	bookmarks := newCollection()

	added := func(key string, value interface{}) {
		bookmark, _ := value.(*Bookmark)
		for _, tag := range bookmark.Tags {
			globalSummary[tag] += 1
		}
	}

	changed := func(key string, value interface{}, old interface{}) {
		addedBooks, _ := value.(*Bookmark)
		removedBooks, _ := old.(*Bookmark)

		for _, tag := range addedBooks.Tags {
			globalSummary[tag] += 1
		}
		for _, tag := range removedBooks.Tags {
			globalSummary[tag] -= 1
		}
	}

	removed := func(key string, value interface{}) {
		bookmark, _ := value.(*Bookmark)
		for _, tag := range bookmark.Tags {
			globalSummary[tag] -= 1
		}
	}

	callback := &Callback{added: added, changed: changed, removed: removed}
	bookmarks.ObserveChanges(callback)

	return &userInfo{subscriptions: newCollection(), bookmarks: bookmarks, summary: make(map[string]int)}
}

func makeUuid() string {
	u, _ := uuid.NewV4()
	return u.String()
}

func hashToken(uuid string) string {
	mac := hmac.New(sha256.New, secretKey)
	mac.Write([]byte(uuid))
	b := mac.Sum(nil)
	return base32.StdEncoding.EncodeToString(b)
}

func init() {

	var err error
	cj, err := ioutil.ReadFile("secret/clientsecrets")

	if err != nil {
		fmt.Println("error reading clientsecrets")
		panic(err.Error())
	}

	clientsecrets := &clientSecrets{}
	err = json.Unmarshal(cj, clientsecrets)

	if err != nil {
		fmt.Println("error unmarshall clientsecret")
		panic(err.Error())
	}

	client_secret = clientsecrets.Client_secret
	client_id = clientsecrets.Client_id
	/*
		secretKey, err = ioutil.ReadFile("secret/secret.key")
		if err != nil {
			panic("i need the secret key")
		}

		//somehow account for trailing newline
		if secretKey[len(secretKey)-1] == 10 {
			secretKey = secretKey[:len(secretKey)-1]
		}
		fmt.Println(secretKey)
	*/
}

func main() {

	http.HandleFunc("/api/summary", summaryHandler)
	http.HandleFunc("/api/img/", authed(imgHandler))
	http.HandleFunc("/ws", authed(websocketHandler))
	http.HandleFunc("/api/follow", authed(subscriptionHandler))
	http.HandleFunc("/api/bookmarks", authed(bookmarkHandler))
	http.HandleFunc("/api/user", userHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	err := http.ListenAndServe(":555", nil)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func authed(h func(w http.ResponseWriter, r *http.Request, context *RequestContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		auth := r.Header.Get("Authorization")

		if auth == "" {
			qs := r.URL.Query()
			auth = qs.Get("user")
		}

		context := &RequestContext{}
		context.isAuthed = true
		context.user = auth

		h(w, r, context)
	}
}
