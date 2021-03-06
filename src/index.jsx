/** @jsx React.DOM */
"use strict";
var HttpClient = require("./httpclient");
var React = require('react');
var ReactDOM = require('react-dom');
var UserSearch = require('./search');
var Bookmarks = require("./bookmarks");

var client = HttpClient();

var storedid = localStorage.getItem("loginid");
var storedtoken =  localStorage.getItem("token");
var accountid = localStorage.getItem("accountid");

function mapQueryString(s){
      return s.split('&')
  .map(function(kvp){
        return kvp.split('=');
  }).reduce(function(sum,current){
        var key = current[0];
        var value = current[1];
        sum[key] = value;
        return sum;
  },{});
}

var qs = mapQueryString(location.search.substr(1));
var code = qs["code"];


var App = React.createClass({
  getInitialState: function() {
    var hostname=location.hostname;
    if (hostname===""){
      hostname="localhost";
    }

    return {
      user:storedid,
      token:storedtoken,
      bookmarks:[],
      subscriptions:[],
      hostname:hostname,
      userlookup:{},
      usernamelookup:{},
      tagfilter:"",
      userid:accountid,
      page:"bookmarks",
      headerpage:"",
      summary: {},
      newbookmark: {Url:"", Description:"", RawTags: ""},
      tagfilters: [],
      tagTimeout: -1,
      isFilterFocused: false,
      isFilterHovered: false
    };
  },
  addSub:function(e){
    var options={
      headers:{Authorization:this.state.user+":"+this.state.token},
      params:{follow:e}
    };

    client.post("/api/follow","",options).then(function(result){
      console.log("followed", result);
    }).catch(function(err){
      console.log(err);
    });

  },
  deleteSub:function(e){
    var options={
      headers:{Authorization:this.state.user+":"+this.state.token},
      params:{follow:e},
      method:"DELETE"
    };

    client.post("/api/follow","",options).then(function(result){
      console.log("followed", result);
    }).catch(function(err){
      console.log(err);
    });

  },
  setPage:function(e,page){

    e.preventDefault();
    this.setState({page:page});
  },
  websocketHandler:function(e){
    var self = this;
    var update = JSON.parse(e.data);
    //console.log("update rec", update);

    if (update.Type === "sub"){
        var id = update.Data;
        var subs = self.state.subscriptions;

        if (update.Op === "add"){
            subs.push(id);
        }

        if (update.Op === "delete"){
            subs = subs.filter(x => x !== id);
        }

        console.log(subs);

        self.setState({subscriptions:subs});
        return
    }

    if (typeof(update.Name)==="undefined"){
      var newbookmarks = self.state.bookmarks.filter(b => b.Id !== update.Id);

      if (update.Url !== ""){
          newbookmarks.push(update);
      }
      //console.log("update:", update.Id, newbookmarks);
      self.setState({bookmarks:newbookmarks});

    }else{
      var usernamelookup = self.state.usernamelookup;
      usernamelookup[update.Userid]=update.Name;
      self.setState({usernamelookup:usernamelookup});
    }
  },
  componentDidMount: function() {
    var self = this;

    client.get("/api/summary").then(function(result){
      var summary = JSON.parse(result);
      self.setState({summary:summary});

    }).catch(err=>console.log(err));

    client.get("/api/user").then(function(result){

      var userlookup = JSON.parse(result);
      var summarylookup = {};
      var namelookup = {};

      Object.keys(userlookup).forEach(function(userid){
        summarylookup[userid]=userlookup[userid].Summary;
        namelookup[userid] = userlookup[userid].Name;
      });

      self.setState({userlookup:summarylookup, usernamelookup:namelookup});
    }).catch(function(err){
      console.log(err);
    });

    var getToken = new Promise(function(resolve,reject){

      if (storedtoken!==null&&storedid!==null&&accountid!==null){
        resolve({user:storedid,token:storedtoken,userid:accountid});
      }else if (typeof(code)!=="undefined") {

        history.pushState('','','/');
        client.post("/api/user?code="+code).then(function(result){
          resolve(JSON.parse(result));
        });

      }else{
          reject("you need to register!");
      }

    });

    getToken.then(function(result){
      self.setState(result);
      localStorage.setItem("loginid",result.user);
      localStorage.setItem("token",result.token);
      localStorage.setItem("accountid",result.userid);

      console.log("logging in as:", result.user+":"+result.token);

      var ws;
      var protocol = location.protocol === "http:" ? "ws:" : "wss:";

      if (location.hostname==="localhost"){
        ws = new WebSocket("ws://"+ self.state.hostname +":555/ws?user="+result.user+"&token="+result.token);
      }else{
        ws = new WebSocket(protocol + "//"+ location.host +"/ws?user="+result.user+"&token="+result.token);
      }
      

      ws.onmessage = self.websocketHandler;

    }).catch(function(err){
      console.error(err);
    });


  },
  filterUsers:function(e){
    var val = e.target.value;
    var last = val.substr(-1);
    var self = this;

    clearTimeout(self.state.tagTimeout);

    var addtag = function(){
      var tag = val.replace(/,/g, "").trim();
      if (tag.length > 0){
        console.log("adding tag:",tag);
        self.addFilter(tag);
        self.setState({tagfilter:""});
      }
    };

    if (last==="," || last===" "){
      addtag();
    }else{
      var timeout = setTimeout(addtag, 1000);
      self.setState({tagfilter:val, tagTimeout:timeout});
    }


  },
  keyDown:function(e){
    var key = e.key;
    if (key==="Backspace" && this.state.tagfilter.length === 0){
      console.log("delete last tag");
      this.popFilter();
    }
  },
  setFilter:function(tag){
    this.setState({tagfilter:tag});
  },
  setHeader:function(e,page){
    e.preventDefault();
    this.setState({headerpage:page});

  },
  addFilter:function(tag){
    var tag = tag.toLowerCase();
    var tags = this.state.tagfilters;

    if (tags.indexOf(tag) >= 0){
      return;
    }

    tags.push(tag);
    this.setState({tagfilters:tags});
  },
  popFilter:function(){
    var tags = this.state.tagfilters;
    tags.pop();
    this.setState({tagfilters:tags});
  },
  removeFilter:function(tag){
    var tags = this.state.tagfilters.filter(x => x !== tag);
    this.setState({tagfilters:tags});
  },
  onFocus:function(e){
    this.setState({isFilterFocused:true});
  },
  onBlur:function(e){
    this.setState({isFilterFocused:false});
  },
  onMouseOver:function(e){
    this.setState({isFilterHovered:true});
  },
  onMouseOut:function(e){
    this.setState({isFilterHovered:false});
  },
  putBookmark:function(bm){

    var tags = bm.RawTags ? bm.RawTags.split(",").map(x=>x.trim()) : [];
    var bookmark = {id: bm.Id, url: bm.Url, tags:tags, description:bm.Description};

    console.log("putting bookmark:", bookmark);

    var client = HttpClient();
    client.request("PUT", "/api/bookmarks?user=" + this.state.user + "&token=" + this.state.token, JSON.stringify(bookmark)).then(function(rep){console.log(rep)}).catch(function(err){console.log(err)});
  },
  postBookmark:function(){
    var bm = this.state.newbookmark;
    var tags = bm.RawTags.split(",").map(x=>x.trim());
    var bookmark = {url: bm.Url, tags: tags, description: bm.Description};

    var client = HttpClient();
    client.post("/api/bookmarks?user=" + this.state.user + "&token=" + this.state.token, JSON.stringify(bookmark)).then(function(rep){console.log(rep)}).catch(function(err){console.log(err)});

  },
  render:function(){
    var self = this;
    var setPage = this.setPage;
    var setHeader = this.setHeader;

    var searchBookmarks = function(e){
      setPage(e,"bookmarks");
    };

    var searchUsers = function(e){
      setPage(e,"users");
    };

    var showBookmarklet = function(e){
      setHeader(e,"bookmarklet");
    };

    var showNewBookmark = function(e){
      setHeader(e, "newbookmark");
    };

    var nothing= function(e){
        setHeader(e,"");
    };

    var summary = self.state.summary;

    var tag_breakdown = Object.keys(summary).map(x=>({tag:x, count:summary[x]}));

    tag_breakdown.sort((a,b)=>b.count-a.count);

    var popular_tags = tag_breakdown.map(x => {
      var onclick = function(e){
        e.preventDefault();
        self.addFilter(x.tag);
      };
      return (<div key={x.tag} className="tag" onClick={onclick}> {x.tag} </div>);
    });

    var currentFilter = this.state.tagfilter.toLowerCase();
    var bookmarks = this.state.bookmarks;
    var currentPage = this.state.page;
    var currentHeader = this.state.headerpage;
    var currentfilters = self.state.tagfilters;
    var subscriptions = self.state.subscriptions;

    var githublogin = "";
    var tevs = "";
    if (this.state.user===null){
      githublogin = <a href="https://github.com/login/oauth/authorize?client_id=f584faa0641263aab644">{"login through github"}</a>
      tevs= <span> | </span>
    }

    bookmarks = bookmarks.filter(function(bookmark){
        var matchingSubs = subscriptions.filter(x=>x == bookmark.Owner);
        return matchingSubs.length>0;
    });

    currentfilters.forEach(function(tag){
      bookmarks = bookmarks.filter(function(bookmark){
        var tags = bookmark.Tags;
        for (var i = 0; i < tags.length; i++){
          if (tags[i].toLowerCase() === tag){
            return true;
          }
        }        
        return false;
      });
    });

    bookmarks.sort((a,b)=>b.Time-a.Time);

    var content = "";

    switch (currentPage){
      case "bookmarks":
      content=(
        <Bookmarks bookmarks={bookmarks} usernamelookup={this.state.usernamelookup} user={this.state.userid} putBookmark={this.putBookmark}></Bookmarks>
      );
      break;
      case "users":
      content= (<UserSearch userlookup={this.state.userlookup} currentFilters={self.state.tagfilters} usernamelookup={this.state.usernamelookup} onsubadded={this.addSub} onsubdeleted = {this.deleteSub}></UserSearch>);
      break;

      default:
      break;
    }

    var part1=
      'javascript:(function(e){typeof(_tevsel)==="undefined"||!(_tevsel&&_tevsel.parentElement&&document.body.removeChild(_tevsel));_tevsel=document.createElement("div");_tevsel.style.position="fixed";_tevsel.style.top="0";_tevsel.style.left="0";_tevsel.style.backgroundColor="cornflowerblue";_tevsel.style.zIndex="9999";_tevsel.style.padding="5px";_tevsurl=document.createElement("input");_tevsurl.type="text";_tevsurl.value=location;_tevsurl.style.width="350px";_tevsurl.style.display="block";_tevsurl.style.margin="3px";_tevsurl.style.padding="3px";_tevsdes=document.createElement("input");_tevsdes.type="text";_tevsdes.style.width="350px";_tevsdes.style.display="block";_tevsdes.style.margin="3px";_tevsdes.style.padding="3px";_tevsdes.placeholder="enter\\x20a\\x20description";_tevsdes.value=document.title;_tevstags=document.createElement("input");_tevstags.type="text";_tevstags.style.width="350px";_tevstags.style.display="block";_tevstags.style.margin="3px";_tevstags.style.padding="3px";_tevstags.placeholder="separate\\x20tags\\x20with\\x20commas";_tevsbut=document.createElement("input");_tevsbut.type="button";_tevsbut.value="submit";_tevsbut.style.padding="8px";_tevsbut.style.margin="3px";_tevscancel=document.createElement("input");_tevscancel.type="button";_tevscancel.value="cancel";_tevscancel.style.margin="3px";_tevscancel.style.padding="8px";_tevscancel.style.float="right";document.body.appendChild(_tevsel);_tevsel.appendChild(_tevsurl);_tevsel.appendChild(_tevsdes);_tevsel.appendChild(_tevstags);_tevsel.appendChild(_tevsbut);_tevsel.appendChild(_tevscancel);_tevscancel.onclick=function(e){document.body.removeChild(_tevsel);};_tevsbut.onclick=function(e){_tevscontent={url:_tevsurl.value,description:_tevsdes.value,tags:_tevstags.value.split(",").map(function(tag){return(tag.trim())})};_tevsclient=document.createElement("img");_tevsclient.src="http://';
      var part3='";document.body.removeChild(_tevsel);}})();';
      var bookmarklet = part1+ this.state.hostname + "/api/img/user.gif?user="+this.state.user+"&token=" + encodeURIComponent(this.state.token) + "&body=\"+encodeURIComponent(JSON.stringify(_tevscontent))+\"" + part3;

    var headerContent = "";

    switch(currentHeader){

      case "bookmarklet":
      headerContent=(<div className="smaller">
                <p><label>your bookmarklet url:<input type="text" value={bookmarklet}></input></label></p>
              </div>);

      break;

      case "newbookmark":

        var setnewbm = function(bm){
          self.setState({newbookmark:bm});
        };

        var changedes = function(e){
          var bm = self.state.newbookmark;
          bm.Description=e.target.value;
          setnewbm(bm);
        };

        var changeurl = function(e){
          var bm = self.state.newbookmark;
          bm.Url = e.target.value;
          setnewbm(bm);
        };

        var changetags = function(e){
          var bm = self.state.newbookmark;
          bm.RawTags = e.target.value;
          setnewbm(bm);
        };

        var submit_bm = function(e){
          e.preventDefault();
          self.postBookmark();
        };

        headerContent = (<div style={{height:"20em"}}>
                  <input type="text" value={self.state.newbookmark.Url} onChange={changeurl} placeholder={"Url"}></input>
                  <input type="text" value={self.state.newbookmark.Description} onChange={changedes} placeholder={"Description"}></input>
                  <input type="text" value={self.state.newbookmark.RawTags} onChange={changetags} placeholder={"Tags"}></input>
                  <button type="button" onClick={submit_bm}>submit!</button>
                  <button type="button" onClick={nothing}>cancel</button>
                 </div>);
      default:

      break;
    }
        
    //<a href="#" onClick={showNewBookmark}>add bookmark</a>

    var tagfilters = self.state.tagfilters.map(function(tag){
      var onclick = function(e){self.removeFilter(tag);};
      return(<span className="tag" key={tag} onClick={onclick}>{tag}  {"\u2716"}</span>);
    });

    var searchColor = "gainsboro";
    var inputStyle = "input";

    if (self.state.isFilterFocused || self.state.isFilterHovered){
      searchColor = "springgreen";
      inputStyle = "input active";
    }

    return(
      <div>
        <div className="section">
          <div style={{textAlign:"center"}}>
        {githublogin} {tevs}
        <a href={bookmarklet}>show bookmarklet / add bookmark</a>
          </div>
        
        {headerContent}
      </div>

        <div className="section">
          <div className={inputStyle} style={{width:"100%", overflow:"hidden", padding:"0"}}>
            <div style={{float:"left", height:"100%", margin:"0", padding:"0"}}>
              {tagfilters}

            </div>
            <div style={{overflow:"hidden"}}>
              <input style={{borderWidth:"0", margin:"0"}} value={self.state.tagfilter} onChange={self.filterUsers} onKeyDown={self.keyDown} onFocus={self.onFocus} onMouseOver={self.onMouseOver} onMouseOut={self.onMouseOut} onBlur={self.onBlur} placeholder="tag filter" type="text"/>
            </div>
         </div>

          filter by topic, or select one of the popular topics below
          
        </div>

        <div className="section">
          {popular_tags}
        </div>

        <div className="section">          
          Search <a href="#" onClick={searchBookmarks}>Bookmarks</a> | <a href="#" onClick={searchUsers}>Users</a>
        </div>

        {githublogin}
        <div className="section">
          {content}
        </div>
      </div>
    )
  }
});

ReactDOM.render(
  <App message="hey"/>,
  document.getElementById("content")
);
