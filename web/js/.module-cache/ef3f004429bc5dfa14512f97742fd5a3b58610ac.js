/** @jsx React.DOM */

var Hello = React.createClass({displayName: 'Hello',
  render: function() {
      return (
        React.DOM.div( {className:"commentBox"}, 
          "Hello, ", this.props.name
        )
      );
    }
});

var StatList = React.createClass({displayName: 'StatList',
  render: function() {
    var statNodes = this.props.data.map(function (stat, index) {
      return React.DOM.div(null, stat);
    });
    return React.DOM.div( {className:"commentList"}, commentNodes);
  }
});

React.renderComponent(Hello({name: 'World'}), 
                      document.getElementById("example"));
