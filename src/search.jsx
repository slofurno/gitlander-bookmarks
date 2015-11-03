/** @jsx React.DOM */
var React = require('react');

module.exports = React.createClass({
  getInitialState: function() {
    return {tagfilter:""};
  },
  componentDidMount:function(){

  },
  filterUsers:function(e){
    this.setState({tagfilter:e.target.value});

  },
  render:function(){

    var onsubadded=this.props.onsubadded;
    var users = this.props.userlookup;
    var tagfilter = this.state.tagfilter.toLowerCase();

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

      var user = users[userid];

      var languages=Object.keys(user);

      var max=0;

      languages.map(function(language){
        return user[language];
      }).forEach(function(cur){
        if (cur>max){
          max=cur;
        }
      });

/*
      languages.forEach(function(language){
        summary = summary + language + "  " + user[language] + "\n";
      });
*/

      var bars = languages.map(function(language){
        var count = user[language];
        var scale= (count/max)*250|0;

        var width=scale+"px";

        var style={
            width:width
        };

        return <div className="bargraph" style={style}>{language}</div>

      });

      var addme = function(e){
        onsubadded(userid);
      };

      return(
      <div className="wrap-bar raised" key={userid}>
        <div className="raised bookmark">
          <button className="raised float-right" type="button" onClick={addme}>subscribe!</button>
          <div style={{"padding":"1em 0"}}>
            user_name@gmail.com
          </div>

          {bars}

        </div>
      </div>)
    });

    return (
     <div>
      <p><label>tag filter: <input onChange={this.filterUsers} type="text"/></label></p>
      {usersummaries}
    </div>)


  }

});
