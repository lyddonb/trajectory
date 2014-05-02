/** @jsx React.DOM */

var Backbone = require('backbone');
Backbone.$ = $;
var React = require('react');

var Machines = require('./stats/machine');
var TaskTree = require('./track');

var Router = Backbone.Router.extend({
  initialize: function() {
    "use strict";

    Backbone.history.start();
  },

  routes: {
    'machines/:path': 'machines',
    'machines*': 'machines',
    'track/:requestId': 'track',
  },

  'machines': function(path) {
    "use strict";

    this.switchReactView(<Machines url={path} pollInterval={2000} />);
  },

  'track': function(requestId) {
    "use strict";

    this.switchReactView(<TaskTree requestId={requestId} />);
  },

  getNode: function() {
    "use strict";

    return document.getElementById("stats");
  },

  unmount: function(node) {
    "use strict";

    React.unmountComponentAtNode(node);
  },

  switchReactView: function(view) {
    "use strict";

    var node = this.getNode();
    this.unmount(node);

    // TODO: Override equals check and check if same view before reloading.
    //var args = Array.prototype.slice.call(arguments);
    //args = args.slice(1, args.length);

    React.renderComponent(view, node);
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
