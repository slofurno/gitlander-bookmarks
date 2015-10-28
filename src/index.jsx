/** @jsx React.DOM */
"use strict";
var HttpClient = require("./httpclient");
var React = require('react');
var ReactDOM = require('react-dom');

var client = HttpClient();

var Welcome = React.createClass({
  getInitialState: function() {
    return {user:"sdfsdf", token:"???"};
  },
  componentDidMount: function() {
    var self = this;

    var hostname=location.hostname;
    if (hostname===""){
      hostname="localhost";
    }

    client.post("/api/user").then(function(result){
      self.setState(result);
      console.log(result);

      var ws = new WebSocket("ws://gitlander.com:555/ws?user="+result.user+"&token="+result.token);

      ws.onmessage=function(e){
        var updates = JSON.parse(e.data);
        console.log(updates);
      };
    });
/*
    console.log(hostname);

    */

  },
  render:function(){

    var part1=
    'javascript:(function(e){typeof(_tevsel)==="undefined"||!(_tevsel&&_tevsel.parentElement&&document.body.removeChild(_tevsel));_tevsel=document.createElement("div");_tevsel.style.position="fixed";_tevsel.style.top="0";_tevsel.style.left="0";_tevsel.style.backgroundColor="cornflowerblue";_tevsel.style.zIndex="9999";_tevsel.style.padding="5px";_tevsurl=document.createElement("input");_tevsurl.type="text";_tevsurl.value=location;_tevsurl.style.width="350px";_tevsurl.style.display="block";_tevsurl.style.margin="3px";_tevsurl.style.padding="3px";_tevsdes=document.createElement("input");_tevsdes.type="text";_tevsdes.style.width="350px";_tevsdes.style.display="block";_tevsdes.style.margin="3px";_tevsdes.style.padding="3px";_tevsdes.placeholder="enter\\x20a\\x20description";_tevstags=document.createElement("input");_tevstags.type="text";_tevstags.style.width="350px";_tevstags.style.display="block";_tevstags.style.margin="3px";_tevstags.style.padding="3px";_tevstags.placeholder="separate\\x20tags\\x20with\\x20commas";_tevsbut=document.createElement("input");_tevsbut.type="button";_tevsbut.value="submit";_tevsbut.style.padding="8px";_tevsbut.style.margin="3px";_tevscancel=document.createElement("input");_tevscancel.type="button";_tevscancel.value="cancel";_tevscancel.style.margin="3px";_tevscancel.style.padding="8px";_tevscancel.style.float="right";document.body.appendChild(_tevsel);_tevsel.appendChild(_tevsurl);_tevsel.appendChild(_tevsdes);_tevsel.appendChild(_tevstags);_tevsel.appendChild(_tevsbut);_tevsel.appendChild(_tevscancel);_tevscancel.onclick=function(e){document.body.removeChild(_tevsel);};_tevsbut.onclick=function(e){_tevscontent={url:_tevsurl.value,description:_tevsdes.value,tags:_tevstags.value.split(",").map(function(tag){return(tag.trim())})};_tevsclient=new(XMLHttpRequest);_tevsclient.onload=function(e){console.log(this.data)};_tevsclient.open("post","http://gitlander.com:555/api/bookmarks");_tevsclient.setRequestHeader("Authorization","'
    var part3='");_tevsclient.send(JSON.stringify(_tevscontent));document.body.removeChild(_tevsel);}})();';

    console.log("RENDERING");

    var fullstring=part1+this.state.user+":"+this.state.token+part3;

    return(
      <div>
        <div className="smaller">
          {fullstring}
        </div>
        <div className="tester raised">

          {this.state.user}
          {this.state.token}
          {this.props.message}
        </div>
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
