/** @jsx React.DOM */

var Machines = React.createClass({displayName: 'Machines',
  loadMachinesFromServer: function() {
    $.ajax({
      url: this.props.url,
      success: function(data) {
        this.setState({data: data});
      }.bind(this)
    });
  },
  componentWillMount: function() {
    this.loadMachinesFromServer();
    // TODO: Clean this up.
    setInterval(this.loadMachinesFromServer, this.props.pollInterval);
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

var MachineList = React.createClass({displayName: 'MachineList',
  render: function() {
    var machineNodes = this.props.data.map(function (machine, index) {
      return machine.split(".").map(function (machineItem, index) {
        return React.DOM.div(null, machineItem);
      });
    });
    return React.DOM.div( {className:"machineList"}, machineNodes);
  }
});

var Stats = React.createClass({displayName: 'Stats',
  loadStatsFromServer: function() {
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

React.renderComponent(Stats( {url:"/api/hosts", pollInterval:2000} ), 
                      document.getElementById("stats"));
