/** @jsx React.DOM */
var React = require('react');

module.exports = React.createClass({
  getInitialState: function() {
    return {};
  },
  componentDidMount:function(){

  },
  render:function(){

    var usernamelookup = this.props.usernamelookup;
    var onsubadded=this.props.onsubadded;
    var users = this.props.userlookup;
    var tagfilter = this.props.currentFilter.toLowerCase();

    console.log(users);

    var allusers = Object.keys(users).filter(function(userid){
      return Object.keys(users[userid]).length>0;
    });

    var filteredusers = allusers;

    if (tagfilter!=""){
      filteredusers=filteredusers.filter(function(userid){
        var languages = Object.keys(users[userid]);
        var matches = languages.filter(function(language){
          return language.substr(0,tagfilter.length).toLowerCase()===tagfilter;
        });
        return matches.length>0;
      });
    }

    var usersummaries = filteredusers.map(function(userid){
      console.log("userid",userid);

      var userName = usernamelookup[userid];
      var user = users[userid];

      var languages=Object.keys(user);

      var max=0;
      var total=0;

      languages.map(function(language){
        return user[language];
      }).forEach(function(cur){
        total+=cur;
        if (cur>max){
          max=cur;
        }
      });

      languages.sort((a,b)=>user[b]-user[a]);
/*
      languages.forEach(function(language){
        summary = summary + language + "  " + user[language] + "\n";
      });
*/

      var bars = languages.map(function(language){
        var count = user[language];
        var scale= (count/max)*100 |0;//*250|0;

        var width=scale+"px";

        var style={
            width:scale+"%"
        };

        return <div className="bargraph" style={style}>{language}</div>

      });

      var addme = function(e){
        onsubadded(userid);
      };

      return(
      <div key={userid} className="bookmark raised" style={{width:"360px", margin:"0 2px 2px 0", padding:"1em"}}>
        <div style={{height:"20em", overflowY:"hidden"}}>
          <button className="raised float-right" type="button" onClick={addme}>subscribe!</button>
          <div style={{"padding":"1em 0"}}>
            {userName}
          </div>

          {bars}
        </div>
      </div>)
    });

    return (
     <div>
      {usersummaries}
    </div>)


  }

});
