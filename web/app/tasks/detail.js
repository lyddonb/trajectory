/** @jsx React.DOM */

var React = require('react');
var Urls = require('../urls');

var TaskDetail = React.createClass({

  componentDidMount: function() {
    $.ajax({
      url: Urls.getTaskKeysForTaskIdUrl(this.props.taskId),
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

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    return <TaskList taskId={this.props.taskId} data={this.state.data} />;
  }
});

var TaskList = React.createClass({

  render: function() {
    var tasks = this.props.data.map(function(taskKey, index) {
      // TODO: Convert to link node.
      return <TaskItem taskKey={taskKey.Key} />
    });

    return <div className="col-md-12">
      <h3>Task {this.props.taskId}</h3>
      {tasks}
    </div>;
  }
});


var TaskItem = React.createClass({

  componentDidMount: function() {
    $.ajax({
      url: Urls.getTaskInfoUrl(this.props.taskKey),
      success: function(data) {
        // TODO: Make this a default handler.
        if (data.success) {
          this.setState({data: data.result});
        } else {
          console.log("Failed to load tasks.")
        }
      }.bind(this)
    });
  },

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    var self = this;

    var tasks = Object.keys(this.state.data).map(function(key, index) {
      return buildFormGroup(key, self.state.data[key]);
    });

    var request = <div></div>;

    if (self.state.data["request_info"] !== undefined) {
      var requestInfoSplit = self.state.data["request_info"].split("#");
      var requestid = requestInfoSplit[requestInfoSplit.length - 1]
      request = <TaskRequestItem requestid={requestid} />
    }

    return <div className="col-md-12">
      <div className="col-md-6">
        <h3>Task Stats</h3>
        {tasks}
      </div>
      {request}
    </div>
  }

});


var TaskRequestItem = React.createClass({

  componentDidMount: function() {
    if (this.props.requestid !== null && this.props.requestid !== undefined &&
                                        this.props.requestid !== "") {
      $.ajax({
        url: Urls.getRequestStatsUrl(this.props.requestid),
        success: function(data) {
          // TODO: Make this a default handler.
          if (data.success) {
            this.setState({data: data.result});
          } else {
            console.log("Failed to load addresses.")
          }
        }.bind(this)
      });
    } else {
      this.setState({data: {}});
    }
  },

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    var self = this;

    var requests = Object.keys(this.state.data).map(function(key, index) {
      return buildFormGroup(key, self.state.data[key]);
    });

    return <div className="col-md-6">
      <h3>Request Stats</h3>
      {requests}
    </div>
  }

});


var buildFormGroup = function(prop, value) {
  if (prop == "ran" || prop == "task_eta") {
    var date = new Date(value*1000);
    var hours = date.getHours();
    var minutes = date.getMinutes();
    var seconds = date.getSeconds();
    value = date.toString();
  }

  return (
    <div className="form-group">
      <label>{prop}:&nbsp;</label>
      <span>{value}</span>
    </div>
  )
}

module.exports = TaskDetail;
