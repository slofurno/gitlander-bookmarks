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

      var user = users[userid];
      var summary="";
      Object.keys(user).forEach(function(language){
        summary = summary + language + "  " + user[language] + "\n";
      });

      var addme = function(e){
        onsubadded(userid);
      };

      return  <div>

        <h3>{userid}</h3>{summary} <button className="raised" type="button" onClick={addme}>subscribe!</button> </div>
    });

    return (<div>
      <p><label>tag filter: <input onChange={this.filterUsers} type="text"/></label></p>
      {usersummaries}
    </div>)


  }

});
