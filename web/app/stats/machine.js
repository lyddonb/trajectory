/** @jsx React.DOM */

var React = require('react');
var Stats = require('./stats');

var MachineHost = "http://localhost:3000/api/tasks/addresses"; 


function getHostUrl() {
  return location.origin + "#machines";
}

String.prototype.endsWith = function(suffix) {
  return this.indexOf(suffix, this.length - suffix.length) !== -1;
};

// TODO: Move to a utils thing or something.
function  buildMachineUrl(url, path) {
  if (path !== null && path !== undefined) {
    url = url + "?path=" + path;
  }

  return url;
}


var Machines = React.createClass({
  loadMachinesFromServer: function() {
    this.props.path = this.props.url;

    var url = buildMachineUrl(MachineHost, this.props.url)

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
    var urlSplit = [];

    if (this.props.url !== undefined && this.props.url !== null) {
      urlSplit = this.props.url.split(".");
    }

    return (
      <div>
        <MachineBreadcrumb data={urlSplit} />
        <MachineList data={this.state.data} />
        <Stats />
      </div>
    );
  }
});

var MachineBreadcrumb = React.createClass({
  render: function() {
    var linkUrl = "";
    var crumbs = this.props.data.map(function(name, index) {
      if (linkUrl !== "" && !linkUrl.endsWith(".")) {
        linkUrl += ".";
      }
      linkUrl += name;

      var url = getHostUrl() + "/" + linkUrl;
      return <li><a href={url}>{name}</a></li>;
    });

    return (
      <ol className="breadcrumb">
        <li><a href={getHostUrl()}>Home</a></li>
        {crumbs}
      </ol>
    )
  }
});

var MachineList = React.createClass({
  render: function() {
    console.log(Object.keys(this.props.data));
    var machineNodes = Object.keys(this.props.data).map(function (machine, index) {
      // TODO: Convert to link node.
      var url = "#/machines/";
      var full_name = machine;

      //if (machine.parent !== null && machine.parent !== undefined) {
        //full_name = machine.parent + "." + machine.machine;
      //}

      url += full_name;

      return <div><a href={url}>{machine}</a></div>;
    });

    return <div className="machineList col-md-2">
      <h3>Servers</h3>
      {machineNodes}
    </div>;
  }
});


module.exports = Machines;
