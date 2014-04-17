var Hello = React.createClass({
  render: function() {
      return React.DOM.div({}, 'Hello ' + this.props.name);
    }
});

React.renderComponent(Hello({name: 'World'}), 
                      document.getElementById("example"));
