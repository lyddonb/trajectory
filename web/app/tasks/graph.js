/** @jsx React.DOM */

var d3 = require('d3');
var React = require('react');

function getTaskRequestGraphUrl(address, request) {
  return "http://localhost:3000/api/tasks/addresses/" + address + "/requests/" +
  request + "/taskgraph";
}

var TaskRequestGraph = React.createClass({

  componentDidMount: function () {
    track(this.props.host, this.props.requestid);
  },

  render: function() {
    "use strict";

    return (
      <div id="taskContainer">
        <div id="taskTree"></div>
      </div>
    );
  }

});


var buildFormGroup = function(prop, value) {
  return (
    <div className="form-group">
      <label>{prop}:&nbsp;</label>
      <span>{value}</span>
    </div>
  )
};


var TaskPopUpList = React.createClass({

  getInitialState: function() {
    return {data: []};
  },

  componentDidMount: function() {
    console.log(this.props);
    $.ajax({
      url: getTaskRequestGraphUrl(this.props.host, this.props.requestid),
      success: function(data) {
        if (data.success) {
          this.setState({data: data.result}).bind(this);
        } else {
          console.log("Failed to load request graph.")
        }
      }
    });
  },

  render: function() {
    var tasks = [],
        counter = 0;

    for (var prop in this.state.data) {
      if (prop == "request_id") { continue; }

      var value = "";

      if (prop != "task_id") {
        value = this.state.data[prop];
      } else {
        var propSplit = this.state.data[prop].split("|");
        prop = "context id";

        if (propSplit.length > 1) {
          value = propSplit[1];
        }
      }

      tasks[counter] = buildFormGroup(prop, value);
      counter++;
    }
    return <div className="taskInfo panel-body">{tasks}</div>;
  }

});


var TaskPopUpLists = React.createClass({

  render: function() {
    //ParentTaskId:TaskId:RequestId
    var id = "taskPane" + this.props.taskId;

    var tasks = this.props.keys.map(function(key, index) {
      var keySplit = key.split(':'),
          show = index ? "" : "in"

      return (
        <div class="panel panel-default">
          <div class="panel-heading">
            <h4 class="panel-title">
              <a data-toggle="collapse" data-parent={"#" + id} href={"#" + keySplit[2]}>
                {keySplit[2]}
              </a>
            </h4>
          </div>
          <div id={keySplit[2]} class="panel-collapse collapse {show}">
            <TaskPopUpList taskKey={key} />
          </div>
        </div>
      )
    });

    return <div id={id} className="popover-content panel-group">{tasks}</div>;
  }

});


var TaskPopUp = React.createClass({

  render: function() {
    "use strict";

    var style = {
      'top': this.props.y + 'px',
      'left': this.props.x + 'px',
      'display': 'inline',
      'max-width': '600px',
      'min-width': '300px'
    }

    console.log(this.props.taskNodeData);

    return (
      <div className="taskTable popover right" style={style}>
        <h3 className="popover-title">
          <b>Task Id:</b> {this.props.taskNodeData.task_id}
        </h3>
        <TaskPopUpLists keys={this.props.taskNodeData.keys} 
          taskId={this.props.taskNodeData.task_id} />
      </div>
    )
  }

});

var track = function(address, requestid) {
  "use strict";

  var m = [10, 10, 10, 10],
    w = 2280 - m[1] - m[3],
    h = 1800 - m[0] - m[2],
    i = 0,
    tree = d3.layout.tree().size([h, w]),
    diagonal = d3.svg.diagonal().projection(function(d) { return [d.x, d.y]; }),

    vis = d3.select("#taskTree").append("svg:svg")
      .attr("width", w + m[1] + m[3])
      .attr("height", h + m[0] + m[2])
    .append("svg:g")
      .attr("transform", "translate(" + (m[3]) + "," + m[0] + ")");

  function toggleAll(d) {
    if (d.children) {
      d.children.forEach(toggleAll);
      toggle(d);
    }
  }

  $.ajax({
    // TODO: Pass in the url.
    url: getTaskRequestGraphUrl(address, requestid),
    success: function(data) {

      if (data.success) {
        var root = data.result;

        root.x0 = h / 2;
        root.y0 = 0;

        update(root, root);

        if (root.children !== undefined && root.children != null) {
          root.children.forEach(toggleAll);
        }

      } else {
        console.log("Failed to load request graph.")
      }

    }.bind(this)
  });

  function show(taskNodeId, taskNodeData) {
    var taskPopUp = $("#" + taskNodeData.task_id);
    if (taskPopUp.length !== 0) {
      taskPopUp.remove();
    } else {
      var taskNode = d3.select(taskNodeId)[0][0].parentNode;
      var selTaskNode = d3.select(taskNode);
      var boundingRect = taskNode.getBoundingClientRect();
      var position = d3.transform(selTaskNode.attr("transform")).translate;
      var x0 = 70;
      var y0 = 70;
      if (d3.select(
        selTaskNode[0][0].childNodes[0]).attr("text-anchor") !== "start") {
        x0 += 60;
      }

      $("#taskContainer").append(
        $('<div></div>').attr("id", taskNodeData.task_id)
      );
      React.renderComponent(
        <TaskPopUp x={position[0] + x0} y={position[1] + y0}
          taskNodeData={taskNodeData} />,
        document.getElementById(taskNodeData.task_id));
    }
  }

  function update(root, source) {
    var duration = d3.event && d3.event.altKey ? 5000 : 500;

    // Compute the new tree layout.
    var nodes = tree.nodes(root);

    // Normalize for fixed-depth.
    nodes.forEach(function(d) { 
      d.y = d.depth * 30; 

      if (d.parent !== undefined && d.parent !== null) {
        if (d.parent.children.length > 1) {
          d.y += d.parent.children.indexOf(d) * 15;
        }
      }
    });

    // Update the nodes…
    var node = vis.selectAll("g.node")
        .data(nodes, function(d) { return d.id || (d.id = ++i); });

    // Enter any new nodes at the parent's previous position.
    var nodeEnter = node.enter().append("svg:g")
        //.attr("id", function(d) { return d.id; })
        .attr("class", "node")
        .attr("transform", function(d) { return "translate(" + source.x0 + "," + source.y0 + ")"; });

    nodeEnter.append("svg:text")
        .attr("x", function(d) { return d.children || d._children ? -10 : 10; })
        .attr("dy", ".35em")
        .attr("text-anchor", function(d) { return d.children || d._children ? "end" : "start"; })
        .text(function(d) { return d.name; })
        .style("fill-opacity", 1e-6);

    nodeEnter.append("svg:title")
      .text(function(d) { return d.name + " - " + d.key; });

    nodeEnter.append("svg:circle")
        .attr("r", 1e-6)
        .style("fill", function(d) { return d._children ? "lightsteelblue" : "#fff"; })
        .on("click", function(d) { show(this, d); });

    // Transition nodes to their new position.
    var nodeUpdate = node.transition()
        .duration(duration)
        .attr("transform", function(d) { return "translate(" + d.x + "," + d.y + ")"; });

    nodeUpdate.select("circle")
        .attr("r", 4.5)
        .style("fill", function(d) { return d._children ? "lightsteelblue" : "#fff"; });

    nodeUpdate.select("text")
        .style("fill-opacity", 1);

    // Transition exiting nodes to the parent's new position.
    var nodeExit = node.exit().transition()
        .duration(duration)
        .attr("transform", function(d) { return "translate(" + source.x + "," + source.y + ")"; })
        .remove();

    nodeExit.select("circle")
        .attr("r", 1e-6);

    nodeExit.select("text")
        .style("fill-opacity", 1e-6);

    // Update the links…
    var link = vis.selectAll("path.link")
        .data(tree.links(nodes), function(d) { return d.target.id; });

    // Enter any new links at the parent's previous position.
    link.enter().insert("svg:path", "g")
        .attr("class", "link")
        .attr("d", function(d) {
          var o = {x: source.x0, y: source.y0};
          return diagonal({source: o, target: o});
        })
      .transition()
        .duration(duration)
        .attr("d", diagonal);

    // Transition links to their new position.
    link.transition()
        .duration(duration)
        .attr("d", diagonal);

    // Transition exiting nodes to the parent's new position.
    link.exit().transition()
        .duration(duration)
        .attr("d", function(d) {
          var o = {x: source.x, y: source.y};
          return diagonal({source: o, target: o});
        })
        .remove();

    // Stash the old positions for transition.
    nodes.forEach(function(d) {
      d.x0 = d.x;
      d.y0 = d.y;
    });
  }

  // Toggle children.
  function toggle(d) {
    if (d.children) {
      d._children = d.children;
      d.children = null;
    } else {
      d.children = d._children;
      d._children = null;
    }
  }

};


module.exports = TaskRequestGraph;
