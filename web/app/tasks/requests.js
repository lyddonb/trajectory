/** @jsx React.DOM */

var React = require('react');

// TODO: Get a decent config.
var BaseTaskAddressUrl = "http://localhost:3000/api/tasks/addresses/";
var RequestUrl = "/requests"

function getHostUrl() {
  return location.origin + "#/tasks";
}

var HostRequests = React.createClass({

  loadHostRequestsFromServer: function() {
    url = BaseTaskAddressUrl + this.props.host + RequestUrl;

    $.ajax({
      url: url,
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
    this.loadHostRequestsFromServer();
  },

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    var urlSplit = [];

    if (this.props.url !== undefined && this.props.url !== null) {
      urlSplit = this.props.url.split(".");
    }

    return <HostRequestList data={this.state.data} host={this.props.host} />;
  }
});

var HostRequestList = React.createClass({
  render: function() {
    var host = this.props.host;
    var hostNodes = Object.keys(this.props.data).map(function(requestid, index) {
      // TODO: Convert to link node.
      var url = "#/tasks/" + host + "/request/" + requestid + "/graph";

      return <div><a href={url}>{requestid}</a></div>;
    });

    return <div className="machineList col-md-2">
      <h3>Requests</h3>
      {hostNodes}
    </div>;
  }
});


module.exports = HostRequests;
