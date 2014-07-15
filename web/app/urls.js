var PORT = "3000"

var Urls = new function() {
  "use strict";

  this.host = location.origin;
  this.api_host = location.origin;
  this.api = this.api_host + "/api";
  this.web = this.host + "/#";

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

  this.getTaskAddressUrl = function() {
    return this.api + "/tasks/addresses";
  };

  this.getTaskHostUrl = function(host) {
    return this.api + "/tasks/addresses/" + host + "/requests";
  };

}();


module.exports = Urls;
