/** @jsx React.DOM */

var React = require('react');

var TaskAddressUrl = "http://localhost:3000/api/tasks/addresses"; 

var Hosts = React.createClass({
  loadHostsFromServer: function() {
    $.ajax({
      url: TaskAddressUrl,
      success: function(data) {
        // TODO: Make this a default handler.
        if (data.success) {
          this.setState({data: data.result});
        } else {
          console.log("Failed to load addresses.")
        }
      }.bind(this)
    });
  },

  componentDidMount: function() {
    this.loadHostsFromServer();
    //this.interval = setInterval(
      //this.loadStatsFromServer, this.props.pollInterval);
  },

  componentWillUnmount: function() {
    //clearInterval(this.interval);
  },

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    var urlSplit = [];

    if (this.props.url !== undefined && this.props.url !== null) {
      urlSplit = this.props.url.split(".");
    }

    return <HostsList data={this.state.data} />;
  }
});

var HostsList = React.createClass({
  render: function() {
    var hostNodes = Object.keys(this.props.data).map(function (host, index) {
      // TODO: Convert to link node.
      var url = "#/tasks/" + host;

      return <div><a href={url}>{host}</a></div>;
    });

    return <div className="machineList col-md-2">
      <h3>Servers</h3>
      {hostNodes}
    </div>;
  }
});


module.exports = Hosts;
