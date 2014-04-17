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

React.renderComponent(Hello({name: 'World'}), 
                      document.getElementById("example"));
