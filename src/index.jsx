/** @jsx React.DOM */
"use strict";

var React = require('react');
var ReactDOM = require('react-dom');

var Welcome = React.createClass({
  getInitialState: function() {
    return {name:"steve"};
  },
  componentDidMount: function() {
    //this.loadCommentsFromServer();
    //setInterval(this.refreshTempests, 5000);
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
