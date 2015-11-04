package main

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"github.com/nu7hatch/gouuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var userTokens = map[string]string{}
var userInfos = map[string]*userInfo{}

var secretKey []byte
var database = newFilebase("test.db")
var dataStore = &DataStore{}

var client_id string
var client_secret string

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
	token         string
}

type RequestContext struct {
	isAuthed bool
	userinfo *userInfo
}

func newUserInfo() *userInfo {
	return &userInfo{subscriptions: newCollection(), bookmarks: newCollection(), summary: make(map[string]int)}
}

func makeUuid() string {
	u, _ := uuid.NewV4()
	return u.String()
}

func init() {
	var err error
	cj, err := ioutil.ReadFile("secret/clientsecrets")

	if err != nil {
		panic(err.Error())
	}

	clientsecrets := &clientSecrets{}
	err = json.Unmarshal(cj, clientsecrets)

	if err != nil {
		fmt.Println(err.Error())
	}

	client_secret = clientsecrets.Client_secret
	client_id = clientsecrets.Client_id

	secretKey, err = ioutil.ReadFile("secret/secret.key")
	if err != nil {
		panic("i need the secret key")
	}

	scanner := bufio.NewScanner(database.Fd)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		du := &DataUnion{}
		b := scanner.Bytes()
		err := json.Unmarshal(b, du)

		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("bad json:", string(b))
			continue
		}

		var userinfo *userInfo
		var ok bool
		userid := du.UserId

		if userinfo, ok = userInfos[userid]; !ok {
			userinfo = newUserInfo()
			userinfo.userid = du.UserId
			fmt.Println("regen user", userinfo.userid)
			userinfo.subscriptions.Add(userid, userid)
			userInfos[userid] = userinfo
		}

		if du.Bookmark != nil {
			dataStore.AddBookmark(userinfo, du.Bookmark)
		} else if du.Sub != "" {
			dataStore.AddSubscription(userinfo, du.Sub)
		} else {
			userTokens[du.Token] = du.UserId
			userinfo.name = du.Name
		}

	}

	//TODO: idk; by reusing the same datastore functions, we were rewriting everything we read
	database.Pls = true

}

func main() {

	http.HandleFunc("/api/img/", authed(imgHandler))
	http.HandleFunc("/ws", authed(websocketHandler))
	http.HandleFunc("/api/follow", authed(subscriptionHandler))
	http.HandleFunc("/api/bookmarks", authed(bookmarkHandler))
	http.HandleFunc("/api/user", userHandler)
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.ListenAndServe(":555", nil)
}

func authed(h func(w http.ResponseWriter, r *http.Request, context *RequestContext)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization")
		context := &RequestContext{}
		auth := r.Header.Get("Authorization")

		authParams := strings.Split(auth, ":")

		var userid string = ""
		var token string = ""

		if len(authParams) == 2 {
			userid = authParams[0]
			token = authParams[1]
		} else {
			qs := r.URL.Query()
			userid = qs.Get("user")
			token = qs.Get("token")
		}

		if userid != "" && token != "" {

			var githubid string
			var user *userInfo
			var ok bool

			if githubid, ok = userTokens[userid]; ok {

				tb, _ := base32.StdEncoding.DecodeString(token)
				mac := hmac.New(sha256.New, secretKey)
				mac.Write([]byte(userid))
				expected := mac.Sum(nil)

				if hmac.Equal(expected, tb) {

					if user, ok = userInfos[githubid]; ok {

						context.isAuthed = true
						context.userinfo = user
						context.userinfo.userid = githubid
						fmt.Println("user authed as: ", githubid)

					}
				}
			}
		}

		h(w, r, context)
	}
}

//TODO: this entire handler is a duplicate, as a workaround for restrictions on mixed content
func imgHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	body := r.URL.Query().Get("body")

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	br := &BookmarkRequest{}
	err := json.Unmarshal([]byte(body), br)

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	buf := []byte(br.Url)
	//resp, err := http.NewRequest("POST", "localhost:8765", bytes.NewBuffer(buf))

	resp, err := http.Post("http://localhost:8765", "text/plain", bytes.NewBuffer(buf))

	defer resp.Body.Close()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	content, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	bookmark := &Bookmark{
		Id:          makeUuid(),
		Url:         br.Url,
		Description: br.Description,
		Tags:        br.Tags,
		Owner:       context.userinfo.userid,
		Time:        getCurrentTime(),
		Summary:     string(content),
	}

	dataStore.AddBookmark(context.userinfo, bookmark)
	fmt.Fprintln(w, bookmark.Id)
}

func websocketHandler(w http.ResponseWriter, req *http.Request, context *RequestContext) {

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ws := upgrade(w, req)

	connection := newUserConnection(context.userinfo.subscriptions, ws)

	func() {
		for {
			read, code, err := ws.Read()
			if err != nil || code == Close {
				return
			}
			fmt.Println(read)
		}
	}()

	connection.onstop()

	fmt.Println("disconnected")
}

func subscriptionHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sub := r.URL.Query().Get("follow")

	var ok bool

	if _, ok = userInfos[sub]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("cant find sub"))
		return
	}

	dataStore.AddSubscription(context.userinfo, sub)
}

func userHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":

		code := r.URL.Query().Get("code")

		if code == "" {
			return
		}

		response, err := http.Post("https://github.com/login/oauth/access_token?client_id="+client_id+"&client_secret="+client_secret+"&code="+code, "application/json", nil)
		defer response.Body.Close()

		fmt.Println("github status code:", response.StatusCode)

		bb, _ := ioutil.ReadAll(response.Body)

		githubqs, err := url.ParseQuery(string(bb))

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		githubtoken := githubqs.Get("access_token")
		githubresponse, err := http.Get("https://api.github.com/user?access_token=" + githubtoken)

		ghb, _ := ioutil.ReadAll(githubresponse.Body)

		github_user := &GithubUserResponse{}
		err = json.Unmarshal(ghb, github_user)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		//we can get here with all failures from github...
		if github_user.Id == 0 {
			fmt.Println("github userid of 0...")
			return
		}

		fmt.Println("github user:", github_user)

		userid := strconv.FormatUint(github_user.Id, 10)
		usertoken := makeUuid()
		userinfo := newUserInfo()
		userinfo.token = usertoken
		userinfo.userid = userid
		userinfo.name = github_user.Login

		dataStore.AddUser(userinfo)

		mac := hmac.New(sha256.New, secretKey)

		mac.Write([]byte(usertoken))
		b := mac.Sum(nil)

		token := base32.StdEncoding.EncodeToString(b)

		userResponse := struct {
			User   string `json:"user"`
			Token  string `json:"token"`
			UserId string `json:"userid"`
		}{usertoken, token, userid}

		jb, err := json.Marshal(userResponse)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(jb)

	case "GET":

		usersummaries := map[string]map[string]int{}
		for key, userinfo := range userInfos {
			usersummaries[key] = userinfo.summary
		}

		j, _ := json.Marshal(usersummaries)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}

}

func bookmarkHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	method := r.Method

	switch method {

	case "POST":
		b, _ := ioutil.ReadAll(r.Body)

		br := &BookmarkRequest{}
		err := json.Unmarshal(b, br)

		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		if !context.isAuthed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		bookmark := &Bookmark{
			Id:          makeUuid(),
			Url:         br.Url,
			Description: br.Description,
			Tags:        br.Tags,
			Owner:       context.userinfo.userid,
		}

		dataStore.AddBookmark(context.userinfo, bookmark)
		fmt.Fprintln(w, bookmark.Id)

	case "GET":
		userid := r.URL.Query().Get("id")

		if userid == "" {
			return
		}

		if user, ok := userInfos[userid]; ok {

			bookmarks := user.bookmarks.Fetch()
			j, _ := json.Marshal(bookmarks)
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
			return
		}

	}

}
