/** @jsx React.DOM */
"use strict";
var HttpClient = require("./httpclient");
var React = require('react');
var ReactDOM = require('react-dom');
var UserSearch = require('./search');

var client = HttpClient();

var storedid = localStorage.getItem("userid");
var storedtoken = localStorage.getItem("token");

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
      userlookup:{}
    };
  },
  addSub:function(e){
    console.log(e);

    var options={
      headers:{Authorization:this.state.user+":"+this.state.token},
      params:{follow:e}
    };

    client.post("/api/follow","",options).then(function(result){
      console.log(result);
    }).catch(function(err){
      console.log(err);
    });

  },
  componentDidMount: function() {
    console.log("componentDidMount");
    var self = this;

    client.get("/api/user").then(function(result){
      self.setState({userlookup:JSON.parse(result)});
    }).catch(function(err){
      console.log(err);
    });

    var getToken;

    if (storedtoken===null||storedid===null){
      getToken = client.post("/api/user").then(function(result){
        return JSON.parse(result);
      });
    }else{
      getToken = Promise.resolve({user:storedid,token:storedtoken});
    }

    getToken.then(function(result){
      self.setState(result);
      localStorage.setItem("userid",result.user);
      localStorage.setItem("token",result.token);

      console.log("logging in as:", result.user+":"+result.token);
      var ws = new WebSocket("ws://"+ self.state.hostname +":555/ws?user="+result.user+"&token="+result.token);
      ws.onmessage=function(e){
        var update = JSON.parse(e.data);
        var newbookmarks=self.state.bookmarks.slice();

        console.log(update);
        newbookmarks.push(update);

        //TODO:
        /*
        var lookup = self.state.userlookup;
        if (typeof lookup[update.Owner]==="undefined"){
          lookup[update.Owner]={};
        }

        var tags = lookup[update.Owner];

        update.Tags.forEach(function(tag){
          if (typeof tags[tag] ==="undefined"){
            tags[tag]=0;
          }
          tags[tag]+=1;
        });
        */

        //,userlookup:lookup
        self.setState({bookmarks:newbookmarks});
      };

    }).catch(function(err){
      console.error(err);
    });


  },
  render:function(){

    var part1=
    'javascript:(function(e){typeof(_tevsel)==="undefined"||!(_tevsel&&_tevsel.parentElement&&document.body.removeChild(_tevsel));_tevsel=document.createElement("div");_tevsel.style.position="fixed";_tevsel.style.top="0";_tevsel.style.left="0";_tevsel.style.backgroundColor="cornflowerblue";_tevsel.style.zIndex="9999";_tevsel.style.padding="5px";_tevsurl=document.createElement("input");_tevsurl.type="text";_tevsurl.value=location;_tevsurl.style.width="350px";_tevsurl.style.display="block";_tevsurl.style.margin="3px";_tevsurl.style.padding="3px";_tevsdes=document.createElement("input");_tevsdes.type="text";_tevsdes.style.width="350px";_tevsdes.style.display="block";_tevsdes.style.margin="3px";_tevsdes.style.padding="3px";_tevsdes.placeholder="enter\\x20a\\x20description";_tevstags=document.createElement("input");_tevstags.type="text";_tevstags.style.width="350px";_tevstags.style.display="block";_tevstags.style.margin="3px";_tevstags.style.padding="3px";_tevstags.placeholder="separate\\x20tags\\x20with\\x20commas";_tevsbut=document.createElement("input");_tevsbut.type="button";_tevsbut.value="submit";_tevsbut.style.padding="8px";_tevsbut.style.margin="3px";_tevscancel=document.createElement("input");_tevscancel.type="button";_tevscancel.value="cancel";_tevscancel.style.margin="3px";_tevscancel.style.padding="8px";_tevscancel.style.float="right";document.body.appendChild(_tevsel);_tevsel.appendChild(_tevsurl);_tevsel.appendChild(_tevsdes);_tevsel.appendChild(_tevstags);_tevsel.appendChild(_tevsbut);_tevsel.appendChild(_tevscancel);_tevscancel.onclick=function(e){document.body.removeChild(_tevsel);};_tevsbut.onclick=function(e){_tevscontent={url:_tevsurl.value,description:_tevsdes.value,tags:_tevstags.value.split(",").map(function(tag){return(tag.trim())})};_tevsclient=document.createElement("img");_tevsclient.src="http://';
    var part3='";document.body.removeChild(_tevsel);}})();';

    console.log("RENDERING");

    var fullstring=part1+ this.state.hostname + ":555/api/img/user.gif?user="+this.state.user+"&token=" + encodeURIComponent(this.state.token) + "&body=\"+encodeURIComponent(JSON.stringify(_tevscontent))+\"" + part3;

    var bookmarklist = this.state.bookmarks.map(function(bookmark){
      var tags = bookmark.Tags.map(function(tag,index){
        return (<div key={index} className="tag"> {tag} </div>)
      });

      console.log("id: ", bookmark.Id);
      return <div key={bookmark.Id} className="bookmark raised">
        <div>{bookmark.Description}</div>
        <div className="smaller"><a href={bookmark.Url}>{bookmark.Url}</a></div>
        <div> {tags}</div>
      </div>

    });



    return(
      <div>
        <div className="smaller">
          <p><label>your bookmarklet url:<input type="text" value={fullstring}></input></label></p>
        </div>
        <div className="smaller">
          <UserSearch userlookup={this.state.userlookup} onsubadded={this.addSub}></UserSearch>
        </div>

        <div className="section">
          {bookmarklist}
        </div>
      </div>
    )
  }
});

ReactDOM.render(
  <App message="hey"/>,
  document.getElementById("content")
);
