

var Urls = new function() {
  "use strict";

  this.host = "http://localhost:3000";
  this.api = this.host + "/api";
  this.web = this.host + "/#"

  this.getTaskRequestGraphUrl = function(address, request) {
    return this.api + "/tasks/addresses/" + address + "/requests/" +
    request + "/taskgraph";
  };

  this.getRequestStatsUrl = function(request) {
    return this.api + "/stats/" + request;
  };

  this.getTaskInfoUrl = function(taskKey) {
    return this.api + "/tasks/task/" + taskKey;
  };

  this.getTaskKeysForTaskIdUrl = function(taskId) {
    return this.api + "/tasks/tasks/" + taskId;
  };

  this.getTaskDetailPage = function(taskId) {
    return this.web + "/tasks/" + taskId + "/detail";
  };

}();


module.exports = Urls;
