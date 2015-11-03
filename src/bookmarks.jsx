/** @jsx React.DOM */
var React = require('react');


module.exports = React.createClass({
  getInitialState:function(){
    return {};

  },
  componentDidMount:function(){

  },
  render:function(){

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
    /*
    var bookmarklist = this.props.bookmarks.map(function(bookmark){
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
    return <div>{bookmarklist}</div>
    */


  }
});
