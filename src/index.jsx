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
    console.log("getInitialState");
    var hostname=location.hostname;
    if (hostname===""){
      hostname="localhost";
    }

    return {
      user:storedid,
      token:storedtoken,
      bookmarks:[],
      hostname:hostname,
      userlookup:{},
      usernamelookup:{},
      tagfilter:"",
      userid:accountid,
      page:"bookmarks",
      headerpage:""
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
  setPage:function(e,page){

    e.preventDefault();
    this.setState({page:page});
  },
  componentDidMount: function() {
    console.log("componentDidMount");
    var self = this;

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

    var getToken;

    var getToken = new Promise(function(resolve,reject){

      if (storedtoken!==null&&storedid!==null&&accountid!==null){
        resolve({user:storedid,token:storedtoken,userid:accountid});
      }else if (typeof(code)!=="undefined") {

        client.post("/api/user?code="+code).then(function(result){
          resolve(JSON.parse(result));
        });

      }else{
          reject("you need to register!");
      }

    });
/*
    if (typeof(code)!=="undefined"storedtoken===null||storedid===null){
      getToken = client.post("/api/user?code=37ffb2fef37aea6578b9").then(function(result){
        return JSON.parse(result);
      });
    }else{
      getToken = Promise.resolve({user:storedid,token:storedtoken});
    }
    */

    getToken.then(function(result){
      self.setState(result);
      localStorage.setItem("loginid",result.user);
      localStorage.setItem("token",result.token);
      localStorage.setItem("accountid",result.userid);

      console.log("logging in as:", result.user+":"+result.token);
      var ws = new WebSocket("ws://"+ self.state.hostname +":555/ws?user="+result.user+"&token="+result.token);
      ws.onmessage=function(e){
        var update = JSON.parse(e.data);
        console.log("update rec", update);

        if (typeof(update.Name)==="undefined"){
          var newbookmarks=self.state.bookmarks.slice();
          newbookmarks.push(update);
          self.setState({bookmarks:newbookmarks});

        }else{
          var usernamelookup = self.state.usernamelookup;
          usernamelookup[update.Userid]=update.Name;
          self.setState({usernamelookup:usernamelookup});
        }

      };

    }).catch(function(err){
      console.error(err);
    });


  },
  filterUsers:function(e){
    this.setState({tagfilter:e.target.value});

  },
  setHeader:function(e,page){
    e.preventDefault();
    this.setState({headerpage:page});

  },
  render:function(){

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

    var nothing= function(e){
        setHeader(e,"");
    };


    var currentFilter = this.state.tagfilter.toLowerCase();
    var bookmarks = this.state.bookmarks;
    var currentPage = this.state.page;
    var currentHeader = this.state.headerpage;


    var githublogin = "";
    var tevs = "";
    if (this.state.user===null){
      githublogin = <a href="https://github.com/login/oauth/authorize?client_id=f584faa0641263aab644">{"login through github"}</a>
      tevs= <span> | </span>
    }


    if (currentFilter.length>0){
      bookmarks = bookmarks.filter(function(bookmark){
        var matches = bookmark.Tags.map(x=>x.toLowerCase()).filter(x=>x.indexOf(currentFilter)===0);
        return matches.length>0;
      });
    }

    bookmarks.sort((a,b)=>b.Time-a.Time);


    var content = "";

    switch (currentPage){
      case "bookmarks":
      content=(
        <Bookmarks bookmarks={bookmarks} usernamelookup={this.state.usernamelookup} user={this.state.userid}></Bookmarks>
      );
      break;
      case "users":
      content= (<UserSearch userlookup={this.state.userlookup} currentFilter={currentFilter} usernamelookup={this.state.usernamelookup} onsubadded={this.addSub}></UserSearch>);
      break;

      default:
      break;
    }

    var headerContent = "";

    switch(currentHeader){

      case "bookmarklet":
      var part1=
      'javascript:(function(e){typeof(_tevsel)==="undefined"||!(_tevsel&&_tevsel.parentElement&&document.body.removeChild(_tevsel));_tevsel=document.createElement("div");_tevsel.style.position="fixed";_tevsel.style.top="0";_tevsel.style.left="0";_tevsel.style.backgroundColor="cornflowerblue";_tevsel.style.zIndex="9999";_tevsel.style.padding="5px";_tevsurl=document.createElement("input");_tevsurl.type="text";_tevsurl.value=location;_tevsurl.style.width="350px";_tevsurl.style.display="block";_tevsurl.style.margin="3px";_tevsurl.style.padding="3px";_tevsdes=document.createElement("input");_tevsdes.type="text";_tevsdes.style.width="350px";_tevsdes.style.display="block";_tevsdes.style.margin="3px";_tevsdes.style.padding="3px";_tevsdes.placeholder="enter\\x20a\\x20description";_tevsdes.value=document.title;_tevstags=document.createElement("input");_tevstags.type="text";_tevstags.style.width="350px";_tevstags.style.display="block";_tevstags.style.margin="3px";_tevstags.style.padding="3px";_tevstags.placeholder="separate\\x20tags\\x20with\\x20commas";_tevsbut=document.createElement("input");_tevsbut.type="button";_tevsbut.value="submit";_tevsbut.style.padding="8px";_tevsbut.style.margin="3px";_tevscancel=document.createElement("input");_tevscancel.type="button";_tevscancel.value="cancel";_tevscancel.style.margin="3px";_tevscancel.style.padding="8px";_tevscancel.style.float="right";document.body.appendChild(_tevsel);_tevsel.appendChild(_tevsurl);_tevsel.appendChild(_tevsdes);_tevsel.appendChild(_tevstags);_tevsel.appendChild(_tevsbut);_tevsel.appendChild(_tevscancel);_tevscancel.onclick=function(e){document.body.removeChild(_tevsel);};_tevsbut.onclick=function(e){_tevscontent={url:_tevsurl.value,description:_tevsdes.value,tags:_tevstags.value.split(",").map(function(tag){return(tag.trim())})};_tevsclient=document.createElement("img");_tevsclient.src="http://';
      var part3='";document.body.removeChild(_tevsel);}})();';
      var fullstring=part1+ this.state.hostname + ":555/api/img/user.gif?user="+this.state.user+"&token=" + encodeURIComponent(this.state.token) + "&body=\"+encodeURIComponent(JSON.stringify(_tevscontent))+\"" + part3;

      headerContent=(<div className="smaller">
                <p><label>your bookmarklet url:<input type="text" value={fullstring}></input></label></p>
              </div>);

      break;
      default:

      break;
    }


    return(
      <div>
        <div className="section">
          <div style={{textAlign:"center"}}>
        {githublogin} {tevs}
        <a href="#" onClick={nothing}>add bookmark</a><span> | </span>
        <a href="#" onClick={showBookmarklet}>show bookmarklet</a>
          </div>
        {headerContent}
      </div>

        <div className="section">
          <p><label>bookmark filter: <input onChange={this.filterUsers} placeholder="separate tags with commas" type="text"/></label></p>
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
