/** @jsx React.DOM */

var Base = {};

var Stats = React.createClass({displayName: 'Stats',
  loadStatsFromServer: function() {
    console.log(this);
    $.ajax({
      url: this.props.url,
      success: function(data) {
        this.setState({data: data});
      }.bind(this)
    });
  },
  componentWillMount: function() {
    this.loadStatsFromServer();
    setInterval(this.loadStatsFromServer, this.props.pollInterval);
  },
  getInitialState: function() {
    return {data: []};
  },
  render: function() {
      return (
        StatList( {data:this.state.data} )
      );
    }
});

var StatList = React.createClass({displayName: 'StatList',
  render: function() {
    var statNodes = this.props.data.map(function (stat, index) {
      return React.DOM.div(null, stat);
    });
    return React.DOM.div( {className:"statList"}, statNodes);
  }
});

React.renderComponent(Stats( {url:"/api/stats", pollInterval:2000} ), 
                      document.getElementById("stats"));
