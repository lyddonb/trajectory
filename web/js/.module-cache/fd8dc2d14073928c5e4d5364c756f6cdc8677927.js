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
    var commentNodes = this.props.data.map(function (comment, index) {
      return Comment( {key:index, author:comment.author}, comment.text);
    });
    return React.DOM.div( {className:"commentList"}, commentNodes);
  }
});

React.renderComponent(Hello({name: 'World'}), 
                      document.getElementById("example"));
