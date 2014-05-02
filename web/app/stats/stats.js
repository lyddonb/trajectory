/** @jsx React.DOM */
var React = require('react');

var StatUrl = "";

function getReqeustTrackUrl() {
  return location.origin + "#track/";
}
// TODO: Actually filter stats?

var Stats = React.createClass({

  loadMachinesFromServer: function() {
    this.props.url = this.buildUrl(this.props.url)

    $.ajax({
      url: this.props.url,
      success: function(data) {
        this.setState({data: data});
      }.bind(this)
    });
  },

  componentDidMount: function() {
    this.loadMachinesFromServer();
    this.interval = setInterval(
      this.loadStatsFromServer, this.props.pollInterval);
  },

  componentWillUnmount: function() {
    clearInterval(this.interval);
  },

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    return <div className="statList col-md-10">
      <h3>Requests</h3>
      <StatList data={this.state.data}/>
    </div>;
  },

  buildUrl: function(path) {
    var url = "http://localhost:8888/api/stats";

    if (path !== null && path !== undefined) {
      url = url + "?path=" + path;
    }

    return url;
  }

});

var StatList = React.createClass({

  // TODO: Add timestamp to display of each request.
  render: function() {
    var statNodes = this.props.data.map(function(stat, index) {
      return <div><a href={getReqeustTrackUrl() + stat.request_id}>{stat.url}</a></div>
    });

    return <div>{statNodes}</div>;
  }

});


module.exports = Stats;
