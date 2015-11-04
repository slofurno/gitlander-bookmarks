/** @jsx React.DOM */
var React = require('react');


module.exports = React.createClass({
  getInitialState:function(){
    return {};

  },
  componentDidMount:function(){

  },
  render:function(){

    var self = this;

    var bookmarklist = this.props.bookmarks.map(function(bookmark){

      var owner = bookmark.Owner;

      var millis = Date.now() - bookmark.Time;
      var seconds = millis/1000;
      var minutes = seconds/60;
      minutes = minutes|0;

      var agemessage = "";

      if (minutes>90){
        var hours = (minutes/60 + 0.5)|0;
        agemessage = hours + " hours ago";
      }else if (minutes > 0){
        agemessage = minutes + " minutes ago";
      }else{
        agemessage = "just now";
      }

      console.log(self.props.user, owner);
      if (owner===self.props.user){
        owner = "you";
      }

      var tags = bookmark.Tags.map(function(tag,index){
        return (<div key={index} className="tag"> {tag} </div>)
      });

      console.log("id: ", bookmark.Id);
      return <div key={bookmark.Id} className="bookmark raised" style={{width:"360px", margin:"0 2px 2px 0", padding:"1em"}}>
      <div style={{height:"20em", overflowY:"hidden"}}>
        <h2>{bookmark.Description}</h2>
        <div className="smaller"><a href={bookmark.Url}>{bookmark.Url}</a></div>
        <div>{bookmark.Summary}</div>
      </div>
        <div> {tags} {agemessage }<div style={{float:"right",padding:"1em 0 0 0"}}>{owner}</div> </div>

      </div>
    });
    return <div>{bookmarklist}</div>



    /*
    var bookmarks = {};

    this.props.bookmarks.forEach(function(bookmark){
      bookmark.Tags.forEach(function(tag){
        if (typeof(bookmarks[tag])==="undefined"){
          bookmarks[tag]=[];
        }
        bookmarks[tag].push(bookmark);
      });
    });

    var asdf = Object.keys(bookmarks).sort().map(function(tag){
      var list = bookmarks[tag].map(function(bookmark){
        return <div style={{margin:"0 0 0.6em 0"}}><div style={{fontSize:"1.2em", color:"mediumspringgreen"}}>{bookmark.Description}</div>
          <div><a href={bookmark.Url}>{bookmark.Url}</a></div></div>;
      });

      return <div style={{margin:"0 0 1em 0"}}><div className="raised bookmark"><h3>{tag}</h3>{list}</div></div>
    });

    return <div>{asdf}</div>
*/




  }
});
