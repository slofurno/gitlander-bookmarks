/** @jsx React.DOM */
"use strict";
var HttpClient = require("./httpclient");
var React = require('react');
var ReactDOM = require('react-dom');

var client = HttpClient();

var Welcome = React.createClass({
  getInitialState: function() {
    return {userid:"", token:""};
  },
  componentDidMount: function() {

    var hostname=location.hostname;
    if (hostname===""){
      hostname="localhost";
    }

    client.post("/api/user").then(function(result){
      console.log(result);
    });
/*
    console.log(hostname);
    var ws = new WebSocket("ws://"+hostname+"/ws?user=748ddaa2-a558-42e7-61d0-6e0bb4899f37");

    ws.onmessage=function(e){
      var updates = JSON.parse(e.data);
      console.log(updates);
    };
    */

  },
  render:function(){

    console.log("RENDERING");

    return(
      <div className="tester raised">
        {this.props.message}
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
