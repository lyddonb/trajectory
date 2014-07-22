/** @jsx React.DOM */

var React = require('react');
var Urls = require('../urls');

var TaskDetail = React.createClass({

  componentDidMount: function() {
    $.ajax({
      url: Urls.getTaskKeysForTaskIdUrl(this.props.taskId),
      success: function(data) {
        // TODO: Make this a default handler.
        if (data.success) {
          this.setState({data: data.result});
        } else {
          console.log("Failed to load addresses.")
        }
      }.bind(this)
    });
  },

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    return <TaskList taskId={this.props.taskId} data={this.state.data} />;
  }
});

var TaskList = React.createClass({

  render: function() {
    var tasks = this.props.data.map(function(taskKey, index) {
      // TODO: Convert to link node.
      return <TaskItem taskKey={taskKey.Key} />
    });

    return <div className="col-md-12">
      <h3>Task {this.props.taskId}</h3>
      {tasks}
    </div>;
  }
});


var TaskItem = React.createClass({

  componentDidMount: function() {
    $.ajax({
      url: Urls.getTaskInfoUrl(this.props.taskKey),
      success: function(data) {
        // TODO: Make this a default handler.
        if (data.success) {
          this.setState({data: data.result});
        } else {
          console.log("Failed to load tasks.")
        }
      }.bind(this)
    });
  },

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    var self = this;

    var mykeys = ["id", "parent_request_id", "parent_task_id", "request_address",
        "request_info", "task_id", "url"]

    var remainder = ["end", "execution_count", "gae_latency_seconds",
        "ran", "retry_count", "run_time", "status_code", "task_eta"]

    function keyIn(key, list) {
        for (i=0; i < list.length; i++) {
            if (key==list[i]){
                return true;
            }
        }
        return false;
    }

    var tasks = Object.keys(this.state.data).map(function(key, index) {
      if (!keyIn(key, mykeys)) {
        return buildTaskForm(key, self.state.data[key]);
      }
    });

    var task_id = Object.keys(this.state.data).map(function(key, index) {
      if (keyIn(key, mykeys)) {
          return buildTaskForm(key, self.state.data[key]);
      }
    });

    var extra = Object.keys(this.state.data).map(function(key, index) {
      if (!keyIn(key, mykeys) && !keyIn(key, remainder)) {
        return buildTaskForm(key, self.state.data[key]);
      }
    });

    var request = <div></div>;

    if (self.state.data["request_info"] !== undefined) {
      var requestInfoSplit = self.state.data["request_info"].split("#");
      var requestid = requestInfoSplit[requestInfoSplit.length - 1]
      request = <TaskRequestItem requestid={requestid} />
    }

    return <div className="col-lg-12">
      <h3>Task IDs</h3>
      <table className="table table-striped">
        <tbody>
        {task_id}
        </tbody>
      </table>
      <div className="col-md-6">
        <h3>Task Stats</h3>
        <table className="table table-striped">
            <tbody>
            {tasks}
            </tbody>
        </table>
        <h3>Extra Task Stats</h3>
        <table className="table table-striped">
            <tbody>
            {extra}
            </tbody>
        </table>
      </div>
      <div className="col-md-6">
          {request}
      </div>
    </div>
  }

});


var TaskRequestItem = React.createClass({

  componentDidMount: function() {
    if (this.props.requestid !== null && this.props.requestid !== undefined &&
                                        this.props.requestid !== "") {
      $.ajax({
        url: Urls.getRequestStatsUrl(this.props.requestid),
        success: function(data) {
          // TODO: Make this a default handler.
          if (data.success) {
            this.setState({data: data.result});
          } else {
            console.log("Failed to load stats.")
          }
        }.bind(this)
      });
    } else {
      this.setState({data: {}});
    }
  },

  getInitialState: function() {
    return {data: []};
  },

  render: function() {
    var self = this;

    var datakeys = ["datastore_v3_BeginTransaction_count", 
        "datastore_v3_BeginTransaction_duration",
        "datastore_v3_BeginTransaction_offset", "datastore_v3_Commit_count",
        "datastore_v3_Commit_duration", "datastore_v3_Commit_offset",
        "datastore_v3_Get_count", "datastore_v3_Get_duration",
        "datastore_v3_Get_offset", "datastore_v3_Put_count",
        "datastore_v3_Put_duration", "datastore_v3_Put_offset",
        "datastore_v3_AddActions_count", "datastore_v3_AddActions_duration",
        "datastore_v3_AddActions_offset", "datastore_v3_RunQuery_count",
        "datastore_v3_RunQuery_duration", "datastore_v3_RunQuery_offset"]

    var memcachekeys = ["memcache_Delete_count", "memcache_Delete_duration",
        "memcache_Delete_offset", "memcache_Get_count", "memcache_Get_duration",
        "memcache_Get_offset", "memcache_Set_count", "memcache_Set_duration",
        "memcache_Set_offset"]

    var originalkeys = ["cpu_usage", "end_cpu", "end_memory", "exec_time",
        "memory_usage", "overhead", "rpc_total_count", "status_code",
        "system_GetSystemStats_count", "system_GetSystemStats_offset",
        "taskqueue_BulkAdd_count", "taskqueue_BulkAdd_duration", 
        "taskqueue_BulkAdd_offset", "urlfetch_Fetch_count", 
        "urlfetch_Fetch_duration", "urlfetch_Fetch_offset"]
            
    function keyIn(key, list) {
      for (i=0; i < list.length; i++) {
        if (key==list[i]){
          return true;
        }
      }
      return false;
    }

    var requests = Object.keys(this.state.data).map(function(key, index) {
      if (!keyIn(key, datakeys) && !keyIn(key, memcachekeys)) {
        return buildFormGroup(key, self.state.data[key]);
      }
    });

    
    var ds = Object.keys(this.state.data).map(function(key, index) {
      if (keyIn(key, datakeys)) {
        return buildFormGroup(key, self.state.data[key]);
      }
    });
    
    var mc = Object.keys(this.state.data).map(function(key, index) {
      if (keyIn(key, memcachekeys)) {
        return buildFormGroup(key, self.state.data[key]);
      }
    });

    var extra = Object.keys(this.state.data).map(function(key, index) {
      if (!keyIn(key, memcachekeys) && !keyIn(key, datakeys) && !keyIn(key, originalkeys)) {
        return buildFormGroup(key, self.state.data[key]);
      }
    });

    return <div >
      <h3>Request - Datastore Stats</h3>
      <table className="table table-striped">
        <tbody>
        {ds}
        </tbody>
      </table>
      <h3>Request - Memcache Stats</h3>
      <table className="table table-striped">
        <tbody>
        {mc}
        </tbody>
      </table>
      <h3>Request Stats</h3>
      <table className="table table-striped">
        <tbody>
        {requests}
        </tbody>
      </table>
      <h3>Extra Request Stats</h3>
      <table className="table table-striped">
        <tbody>
        {extra}
        </tbody>
      </table>
    </div>
  }

});

var buildTaskForm= function(prop, value) {

    function toTitleCase(str)
    {
        return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() 
            + txt.substr(1).toLowerCase();});
    }

    var mystring = prop;
    mystring = mystring.replace(/_/g, ' ');
    mystring = toTitleCase(mystring);

    var color = " ";

    switch (prop) {
        case "execution_count":
            if (value > 0){
                color = "darkest";
            }
            break;
        case "retry_count":
            if (value > 0){
                color = "darkest";
            }
            break;
        case "status_code":
            if (value != 200){
                color = "darkest";
            }
            break;
        case "gae_latency_seconds":
            if (value > 120){
                color = "darker";
            }else if (value > 60){
                color = "dark";
            }else if (value > 30){
                color = "light";
            }else if (value > 10){
                color = "lighter";
            }else if (value > 1){
                color = "lightest";
            }else{
                color = "";
            }
            break;
        case "run_time":
            if (value > 480){
                color = "darkest";
            }else if (value > 120){
                color = "darker";
            }else if (value > 60){
                color = "dark";
            }else if (value > 10){
                color = "light";
            }else if (value > 5){
                color = "lighter";
            }else if (value > 1){
                color = "lightest";
            }else{
                color = "";
            }
            break;
        default:
            color = " ";
    }

    if (prop == "ran" || prop == "task_eta" || prop == "end") {
      value = new Date(value*1000).toString();
    }

  return (
      <tr className={color}>
          <td title={prop}  className="col-xs-2"><label>{mystring}:&nbsp;</label></td>
          <td className="col-xs-4"><span>{value}</span></td>
      </tr>
  )
}

var buildFormGroup = function(prop, value) {

    function toTitleCase(str)
    {
        return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() 
            + txt.substr(1).toLowerCase();});
    }

    var mystring = prop;
    mystring = mystring.replace(/_/g, ' ');
    mystring = toTitleCase(mystring);

    newValue = value.replace(/\D/g,'');

    var color = " ";

    switch (prop) {
        case "cpu_usage":
            if (newValue > 100){
                color = "darkest";
            }
            break;
        case "memory_usage":
            if (newValue > 200){
                color = "darkest";
            }
            break;
        case "rpc_total_count":
            if (newValue > 100){
                color = "darkest";
            }
            break;
        case "taskqueue_BulkAdd_duration":
            if (newValue > 5){
                color = "darkest";
            }
            break;
        case "status_code":
            if (newValue != 200){
                color = "darkest";
            }
            break;
        case "exec_time":
            if (newValue > 480){
                color = "darkest";
            }else if (newValue > 120){
                color = "darker";
            }else if (newValue > 60){
                color = "dark";
            }else if (newValue > 10){
                color = "light";
            }else if (newValue > 5){
                color = "lighter";
            }else if (newValue > 1){
                color = "lightest";
            }else{
                color = "";
            }
            break;
        case "datastore_v3_Get_duration":
            if (newValue > 20){
                color = "darkest";
            }else if (newValue > 10){
                color = "dark";
            }else if (newValue > 5){
                color = "light";
            }else{
                color = "";
            }
            break;
        case "datastore_v3_Put_duration":
            if (newValue > 20){
                color = "darkest";
            }else if (newValue > 10){
                color = "dark";
            }else if (newValue > 5){
                color = "light";
            }else{
                color = "";
            }
            break;
        case "memcache_Delete_duration":
            if (newValue > 20){
                color = "darkest";
            }else if (newValue > 10){
                color = "dark";
            }else if (newValue > 5){
                color = "light";
            }else{
                color = "";
            }
            break;
        case "memcache_Get_offset":
            if (newValue > 20){
                color = "darkest";
            }else if (newValue > 10){
                color = "dark";
            }else if (newValue > 5){
                color = "light";
            }else{
                color = "";
            }
            break;
        case "memcache_Set_duration":
            if (newValue > 20){
                color = "darkest";
            }else if (newValue > 10){
                color = "dark";
            }else if (newValue > 5){
                color = "light";
            }else{
                color = "";
            }
            break;
        case "urlfetch_Fetch_duration":
            if (newValue > 20){
                color = "darkest";
            }else if (newValue > 10){
                color = "dark";
            }else if (newValue > 5){
                color = "light";
            }else{
                color = "";
            }
            break;
        default:
            color = " ";
    }

  return (
      <tr className={color}>
          <td title={prop} className="col-xs-4"><label>{mystring}:&nbsp;</label></td>
          <td className="col-xs-2"><span>{value}</span></td>
      </tr>
  )
}

module.exports = TaskDetail;
