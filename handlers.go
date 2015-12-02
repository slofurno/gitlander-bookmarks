package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"gitlander.com/slofurno/ws"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

//TODO: this entire handler is a duplicate, as a workaround for restrictions on mixed content
func imgHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	var err error
	body := r.URL.Query().Get("body")

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	br := &BookmarkRequest{}
	err = json.Unmarshal([]byte(body), br)

	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	buf := []byte(br.Url)
	fmt.Println("req to node")
	resp, err := http.Post("http://localhost:8765", "text/plain", bytes.NewBuffer(buf))

	var content []byte

	if err == nil {
		defer resp.Body.Close()
		content, err = ioutil.ReadAll(resp.Body)
	} else {
		fmt.Println("scraper err")
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

	dataStore.UpsertBookmark(context.userinfo, bookmark)
	fmt.Fprintln(w, bookmark.Id)
}

func websocketHandler(w http.ResponseWriter, req *http.Request, context *RequestContext) {

	if !context.isAuthed {
		fmt.Println("not authed")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sock := ws.Upgrade(w, req)

	connection := newUserConnection(context.userinfo.subscriptions, sock)
	defer connection.onstop()

	func() {
		for {
			read, code, err := sock.Read()
			if err != nil || code == ws.Close {
				return
			}
			fmt.Println(read)
		}
	}()

	sock.Close()

	fmt.Println("disconnected")
}

func subscriptionHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	if !context.isAuthed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sub := r.URL.Query().Get("follow")

	var ok bool
	var subinfo *userInfo

	if subinfo, ok = userInfos[sub]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "cant find sub")
		//w.Write([]byte("cant find sub"))
		return
	}

	switch r.Method {
	case "DELETE":
		dataStore.DeleteSubscription(context.userinfo, sub)
	case "POST":
		dataStore.AddSubscription(context.userinfo, sub, subinfo.name)

	}

	w.WriteHeader(http.StatusOK)
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
		userinfo := newUserInfo()
		userinfo.userid = userid
		usertoken := makeUuid()
		userinfo.name = github_user.Login

		dataStore.AddUser(userinfo, usertoken)
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

		usersummaries := map[string]*userinfoDto{}
		for userid, userinfo := range userInfos {
			usersummaries[userid] = &userinfoDto{Summary: userinfo.summary, Name: userinfo.name}
		}

		j, _ := json.Marshal(usersummaries)
		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		return
	}

}

func summaryHandler(w http.ResponseWriter, rep *http.Request) {
	j, _ := json.Marshal(globalSummary)
	fmt.Println(string(j))

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func bookmarkHandler(w http.ResponseWriter, r *http.Request, context *RequestContext) {

	method := r.Method
	var err error

	switch method {

	case "DELETE":
		if !context.isAuthed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		b, _ := ioutil.ReadAll(r.Body)
		br := &BookmarkRequest{}
		err := json.Unmarshal(b, br)

		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		bookmark := &Bookmark{
			Id:          br.Id,
			Url:         br.Url,
			Description: br.Description,
			Tags:        br.Tags,
			Owner:       context.userinfo.userid,
			Time:        getCurrentTime(),
		}

		dataStore.UpsertBookmark(context.userinfo, bookmark)
		w.WriteHeader(http.StatusOK)

	case "PUT":

		if !context.isAuthed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		b, _ := ioutil.ReadAll(r.Body)
		br := &BookmarkRequest{}
		err = json.Unmarshal(b, br)

		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		var content []byte

		if len(br.Url) > 4 {

			buf := []byte(br.Url)
			//TODO: replace with survey w/ timeout?
			resp, err := http.Post("http://localhost:8765", "text/plain", bytes.NewBuffer(buf))

			if err != nil {
				fmt.Println(err.Error())
				return
			}

			//code duplication is bad
			defer resp.Body.Close()
			content, err = ioutil.ReadAll(resp.Body)

			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}

		bookmark := &Bookmark{
			Id:          br.Id,
			Url:         br.Url,
			Description: br.Description,
			Tags:        br.Tags,
			Owner:       context.userinfo.userid,
			Time:        getCurrentTime(),
			Summary:     string(content),
		}

		dataStore.UpsertBookmark(context.userinfo, bookmark)
		w.WriteHeader(http.StatusOK)

	case "POST":

		if !context.isAuthed {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		b, _ := ioutil.ReadAll(r.Body)

		br := &BookmarkRequest{}
		err := json.Unmarshal(b, br)

		buf := []byte(br.Url)
		//TODO: replace with survey w/ timeout?
		resp, err := http.Post("http://localhost:8765", "text/plain", bytes.NewBuffer(buf))

		if err != nil {
			fmt.Println("guess we had an error...")
			//fmt.Println(err.Error())
			return
		}

		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			fmt.Println("guess we had an error reading?")
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

		dataStore.UpsertBookmark(context.userinfo, bookmark)
		fmt.Fprintln(w, bookmark.Id)
	case "GET":
		fmt.Println("when am i calling this?")
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
