var Hello = React.createClass({
  render: function() {
      return (
        <div className="commentBox">
          Hello, {this.props.name}
        </div>
      );
    }
});

/** @jsx React.DOM */
React.renderComponent(Hello({name: 'World'}), 
                      document.getElementById("example"));
