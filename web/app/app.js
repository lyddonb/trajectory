/** @jsx React.DOM */

var Machines = React.createClass({
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
        <MachineList data={this.state.data} />
      );
    }
});

var MachineList = React.createClass({
  render: function() {
    var machineNodes = this.props.data.map(function (machine, index) {
      return machine.split(".").map(function (machineItem, index) {
        return <div>{machineItem}</div>;
      });
    });
    return <div className="machineList">{machineNodes}</div>;
  }
});

var Stats = React.createClass({
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
        <StatList data={this.state.data} />
      );
    }
});

var StatList = React.createClass({
  render: function() {
    var statNodes = this.props.data.map(function (stat, index) {
      return <div>{stat}</div>;
    });
    return <div className="statList">{statNodes}</div>;
  }
});

React.renderComponent(<Machines url="/api/hosts" pollInterval={2000} />, 
                      document.getElementById("stats"));
