/** @jsx React.DOM */

var React = require('react');


function getHostUrl() {
  return location.origin + "#/tasks";
}

String.prototype.endsWith = function(suffix) {
  return this.indexOf(suffix, this.length - suffix.length) !== -1;
};

var Breadcrumb = React.createClass({

  render: function() {
    var crumbs = this.props.data.map(function(crumb, index) {
      var url = getHostUrl() + "/" + crumb.url;
      // TODO: For the last entry don't make it a link?
      return <li><a href={url}>{crumb.name}</a></li>;
    });

    return (
      <ol className="breadcrumb">
        <li><a href={getHostUrl()}>Home</a></li>
        {crumbs}
      </ol>
    )
  }
});

module.exports = Breadcrumb;
