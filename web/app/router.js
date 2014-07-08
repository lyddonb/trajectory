/** @jsx React.DOM */

var Backbone = require('backbone');
Backbone.$ = $;
var React = require('react');

//var Machines = require('./stats/machine');
//var TaskTree = require('./track');
var Breadcrumb = require('./breadcrumb');
var Hosts = require('./tasks/hosts');
var HostRequests = require('./tasks/requests');
var TaskRequestGraph = require('./tasks/graph');
var TaskDetail = require('./tasks/detail');

var Router = Backbone.Router.extend({
  initialize: function() {
    "use strict";

    Backbone.history.start();
  },

  routes: {
    'tasks/:host/request/:requestid/graph': 'taskRequestGraph',
    'tasks/:taskId/detail': 'taskPage',
    'tasks/:host': 'taskRequests',
    'tasks*': 'tasks'
  },

  'tasks': function() {
    this.switchBreadcrumbs(<Breadcrumb data={[]} />);
    this.switchView(<Hosts />);
  },

  'taskRequests': function(host) {
    var crumbs = [
      {url: 'tasks/' + host,  name: host}
    ];

    this.switchBreadcrumbs(<Breadcrumb data={crumbs} />);
    this.switchView(<HostRequests host={host} />);
  },

  'taskRequestGraph': function(host, requestid) {
    var crumbs = [
      {url: 'tasks/' + host,  name: host},
      {url: 'tasks/' + host + "/request/" + requestid + "/graph",  name: requestid}
    ];

    this.switchBreadcrumbs(<Breadcrumb data={crumbs} />);
    this.switchView(<TaskRequestGraph host={host} requestid={requestid} />);
  },

  'taskPage': function(taskId) {
    this.switchView(<TaskDetail taskId={taskId} />);
  },

  //routes: {
    //'machines/:path': 'machines',
    //'machines*': 'machines',
    //'track/:requestId': 'track',
  //},

  //'machines': function(path) {
    //"use strict";

    //this.switchView(<Machines url={path} pollInterval={2000} />);
  //},

  //'track': function(requestId) {
    //"use strict";

    //this.switchView(<TaskTree requestId={requestId} />);
  //},

  getContentNode: function() {
    return document.getElementById("mainContent");
  },

  getBreadcrumbNode: function() {
    return document.getElementById("mainBreadcrumb");
  },

  unmount: function(node) {
    React.unmountComponentAtNode(node);
  },

  switchReactView: function(view, node) {
    this.unmount(node);

    React.renderComponent(view, node);
  },

  switchBreadcrumbs: function(view) {
    this.switchReactView(view, this.getBreadcrumbNode());
  },

  switchView: function(view) {
    this.switchReactView(view, this.getContentNode());
  }
});


// TODO: Make this suck less. Hook into router and bind to the element that will
// be replaced, etc.
$(document).ready(function () {
  "use strict";

  $(document).ajaxSend(function() {
    if ($("#stats > #load").length === 0) {
      $("#stats").append('<i id="load" class="fa fa-spinner fa-spin"></i>');
    }
  });

  $(document).ajaxComplete(function() {
    $("#stats > #load").remove();
  });
});


module.exports = Router;
