/** @jsx React.DOM */
var React = require('react');
var Bookmark = require('./bookmark')


module.exports = React.createClass({
  getInitialState:function(){
    return {};

  },
  componentDidMount:function(){

  },

  render:function(){

    var self = this;
    var user = self.props.user;
    var usernamelookup = self.props.usernamelookup;

    var bookmarklist = this.props.bookmarks.map(function(bookmark){
      //TODO:maybe pass owned bool instead of lookup ref
      return <Bookmark key={bookmark.Id} bookmark={bookmark} user={user} usernamelookup={usernamelookup} putBookmark={self.props.putBookmark}> </Bookmark>

    });

    return <div>{bookmarklist}</div>


      /*
      var owner = bookmark.Owner;

      var millis = Date.now() - bookmark.Time;
      var seconds = millis/1000;
      var minutes = seconds/60;
      minutes = minutes|0;

      var agemessage = "";

      if (minutes>2160){
        var days = (minutes/1440 + 0.5)|0;
        agemessage = days + " days ago";
      }else if (minutes>90){
        var hours = (minutes/60 + 0.5)|0;
        agemessage = hours + " hours ago";
      }else if (minutes > 0){
        agemessage = minutes + " minutes ago";
      }else{
        agemessage = "just now";
      }

      console.log(self.props.user, owner);
      if (owner===self.props.user){
        owner = <a href="#" onClick={startEdit}>Edit</a>
      }else{
        owner = self.props.usernamelookup[owner];
      }

      var tags = bookmark.Tags.map(function(tag,index){
        return (<div key={index} className="tag"> {tag} </div>)
      });

      console.log("id: ", bookmark.Id);
      return (<div key={bookmark.Id} className="bookmark raised" style={{width:"360px", margin:"0 2px 2px 0", padding:"1em"}}>
                <div style={{height:"20em", overflowY:"hidden"}}>
                  <h2>{bookmark.Description}</h2>
                  <div className="smaller"><a href={bookmark.Url}>{bookmark.Url}</a></div>
                  <div>{bookmark.Summary}</div>
                </div>
                <div>
                  {tags} {agemessage }<div style={{float:"right",padding:"1em 0 0 0"}}>{owner}</div>
                </div>

              </div>);
    });

      return <div>{bookmarklist}</div>
    */

  }
});
