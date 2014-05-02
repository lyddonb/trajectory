/** @jsx React.DOM */

//var React = require('react')
var $ = require('jquery');
var Backbone = require('backbone');
Backbone.$ = $;
var Router = require('./router');

//var Stats = React.createClass({
  //loadStatsFromServer: function() {
    //$.ajax({
      //url: this.props.url,
      //success: function(data) {
        //this.setState({data: data});
      //}.bind(this)
    //});
  //},
  //componentWillMount: function() {
    //this.loadStatsFromServer();
    //setInterval(this.loadStatsFromServer, this.props.pollInterval);
  //},
  //getInitialState: function() {
    //return {data: []};
  //},
  //render: function() {
      //return (
        //<StatList data={this.state.data} />
      //);
    //}
//});

//var StatList = React.createClass({
  //render: function() {
    //var statNodes = this.props.data.map(function (stat, index) {
      //return <div>{stat}</div>;
    //});
    //return <div className="statList">{statNodes}</div>;
  //}
//});

if (typeof window !== 'undefined') {
  window.onload = function() {
    "use strict";

    new Router();
  };
}
