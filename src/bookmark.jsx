var React = require('react');

module.exports = React.createClass({
  getInitialState:function(){
    return {isEditing:false};

  },
  componentDidMount:function(){

  },
  saveEdit:function(e){
    e.preventDefault();

    var s = this.state;
    this.props.putBookmark({Id: this.props.bookmark.Id, Url:s.Url, Description:s.Description, RawTags:s.Tags});

    this.replaceState({isEditing:false});
  },
  startEdit:function(e){
    console.log("hi");
    e.preventDefault();
    var bookmark = this.props.bookmark;
    var tags = this.props.bookmark.Tags.reduce((a,c)=>a+=","+c);

    this.setState({isEditing:true, Url:bookmark.Url, Description:bookmark.Description, Tags:tags});
  },
  cancelEdit:function(e){
    e.preventDefault();
    this.replaceState({isEditing:false});
  },
  updateDescription:function(e){
    this.setState({Description:e.target.value});
  },
  updateUrl:function(e){
    this.setState({Url:e.target.value});
  },
  updateTags:function(e){
    this.setState({Tags:e.target.value});
  },
  render:function(){
    var self = this;
    var bookmark = this.props.bookmark;
    var user = this.props.user;
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

    var tags = bookmark.Tags.map(function(tag,index){
      return (<div key={index} className="tag"> {tag} </div>)
    });

    var contents = "";

    if (self.state.isEditing){
      contents = (<div style={{height:"20em", overflowY:"hidden"}}>
                  <input type="text" value={self.state.Description} onChange={self.updateDescription}></input>
                  <input type="text" value={self.state.Url} onChange={self.updateUrl}></input>
                  <input type="text" value={self.state.Tags} onChange={self.updateTags}></input>
                 </div>);
      owner = (<div>
                <a href="#" onClick={self.saveEdit}>Save</a>
                <span> | </span>
                <a href="#" onClick={self.cancelEdit}>Cancel</a>
               </div>
              );

    }else{
      contents = (<div style={{height:"20em", overflowY:"hidden"}}>
                  <h2>{bookmark.Description}</h2>
                 <div className="smaller"><a href={bookmark.Url}>{bookmark.Url}</a></div>
                 <div>{bookmark.Summary}</div>
                 </div>);
      if (owner===self.props.user){
        owner = <a href="#" onClick={self.startEdit}>Edit</a>
      }else{
        owner = self.props.usernamelookup[owner];
      }
    }

    return (<div className="bookmark raised" style={{width:"360px", margin:"0 2px 2px 0", padding:"1em"}}>
             {contents}
              <div>
                {tags} {agemessage}
                <div style={{float:"right",padding:"1em 0 0 0"}}>
                  {owner}
                </div>
              </div>
            </div>);
  }

});
