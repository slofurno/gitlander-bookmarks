/** @jsx React.DOM */
"use strict";
var HttpClient = require("./httpclient");
var React = require('react');
var ReactDOM = require('react-dom');

var client = HttpClient();

var Welcome = React.createClass({
  getInitialState: function() {
    var hostname=location.hostname;
    if (hostname===""){
      hostname="localhost";
    }

    return {
      user:"sdfsdf",
      token:"???",
      bookmarks:[],
      hostname:hostname,
      userlookup:{}
    };
  },
  componentDidMount: function() {
    var self = this;

    client.get("/api/user").then(function(result){
      self.setState({userlookup:result});
    });

    client.post("/api/user").then(function(result){
      self.setState(result);
      console.log(result);

      var ws = new WebSocket("ws://"+ self.state.hostname +":555/ws?user="+result.user+"&token="+result.token);

      ws.onmessage=function(e){
        var update = JSON.parse(e.data);
        var newbookmarks=self.state.bookmarks.slice();
        //var dd=JSON.parse(update.Url)
        //dd.Id=update.Id;
        console.log(update);
        newbookmarks.push(update);

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



/*
        console.log(updates);

        updates.Bookmarks.forEach(function(update){
          var dd=JSON.parse(update.Url);
          //until i fix this, copy the id over
          dd.Id=update.Id;
          newbookmarks.push(dd);
        });
*/
        self.setState({bookmarks:newbookmarks,userlookup:lookup});
      };
    });
/*
    console.log(hostname);

    */

  },
  render:function(){

    var part1=
    'javascript:(function(e){typeof(_tevsel)==="undefined"||!(_tevsel&&_tevsel.parentElement&&document.body.removeChild(_tevsel));_tevsel=document.createElement("div");_tevsel.style.position="fixed";_tevsel.style.top="0";_tevsel.style.left="0";_tevsel.style.backgroundColor="cornflowerblue";_tevsel.style.zIndex="9999";_tevsel.style.padding="5px";_tevsurl=document.createElement("input");_tevsurl.type="text";_tevsurl.value=location;_tevsurl.style.width="350px";_tevsurl.style.display="block";_tevsurl.style.margin="3px";_tevsurl.style.padding="3px";_tevsdes=document.createElement("input");_tevsdes.type="text";_tevsdes.style.width="350px";_tevsdes.style.display="block";_tevsdes.style.margin="3px";_tevsdes.style.padding="3px";_tevsdes.placeholder="enter\\x20a\\x20description";_tevstags=document.createElement("input");_tevstags.type="text";_tevstags.style.width="350px";_tevstags.style.display="block";_tevstags.style.margin="3px";_tevstags.style.padding="3px";_tevstags.placeholder="separate\\x20tags\\x20with\\x20commas";_tevsbut=document.createElement("input");_tevsbut.type="button";_tevsbut.value="submit";_tevsbut.style.padding="8px";_tevsbut.style.margin="3px";_tevscancel=document.createElement("input");_tevscancel.type="button";_tevscancel.value="cancel";_tevscancel.style.margin="3px";_tevscancel.style.padding="8px";_tevscancel.style.float="right";document.body.appendChild(_tevsel);_tevsel.appendChild(_tevsurl);_tevsel.appendChild(_tevsdes);_tevsel.appendChild(_tevstags);_tevsel.appendChild(_tevsbut);_tevsel.appendChild(_tevscancel);_tevscancel.onclick=function(e){document.body.removeChild(_tevsel);};_tevsbut.onclick=function(e){_tevscontent={url:_tevsurl.value,description:_tevsdes.value,tags:_tevstags.value.split(",").map(function(tag){return(tag.trim())})};_tevsclient=document.createElement("img");_tevsclient.src="http://';
    var part3='";document.body.removeChild(_tevsel);}})();';

    console.log("RENDERING");
    //http://gitlander.com:555/api/img/?'
    var fullstring=part1+ this.state.hostname + ":555/api/img/user.gif?user="+this.state.user+"&token=" + encodeURIComponent(this.state.token) + "&body=\"+encodeURIComponent(JSON.stringify(_tevscontent))+\"" + part3;

    var bookmarklist = this.state.bookmarks.map(function(bookmark){
      var tags = bookmark.Tags.map(function(tag,index){
        return (<div key={index} className="tag"> {tag} </div>)
      });

      console.log("id: ", bookmark.Id);
      return <div key={bookmark.Id} className="tester raised">
        <div>{bookmark.Description}</div>
         <div>{bookmark.Url}</div>
        <div> {tags}</div>
      </div>

    });

    var users = this.state.userlookup;

    console.log(users);

    var usersummaries = Object.keys(users).filter(function(userid){
      return Object.keys(users[userid]).length>0;
    }).map(function(userid){

      var user = users[userid];
      var summary="";
      Object.keys(user).forEach(function(language){
        summary = summary + language + "  " + user[language] + "\n";
      });

      return  <div><h3>{userid}</h3> {summary}</div>
    });

    return(
      <div>
        <div className="smaller">
          {fullstring}
        </div>
        <div className="smaller">
          {usersummaries}
        </div>

        {bookmarklist}
      </div>
    )
  }
});

ReactDOM.render(
  <Welcome
  message="HEY"
  />,
  document.getElementById("content")
);
