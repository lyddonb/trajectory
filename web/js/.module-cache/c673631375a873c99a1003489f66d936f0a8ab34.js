/** @jsx React.DOM */

var Hello = React.createClass({displayName: 'Hello',
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

React.renderComponent(Hello({name: 'World'}), 
                      document.getElementById("example"));
