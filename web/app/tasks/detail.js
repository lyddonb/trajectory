/** @jsx React.DOM */

var React = require('react');
var Urls = require('../urls');

var lightest = "lightest";
var lighter = "lighter";
var light = "light";
var dark = "dark";
var darker = "darker";
var darkest = "darkest";

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

    var mykeys = ["id", "parent_task_id", "task_id", "url", "context_id"];

    var remainder = ["end", "execution_count", "gae_latency_seconds",
        "ran", "retry_count", "run_time", "status_code", "task_eta"];

    var request_info = ["app_id", "host", "instance_id", "module_id",
      "request_id", "version_id", "parent_request_id",
      "request_address"];

    function keyIn(key, list) {
        for (i=0; i < list.length; i++) {
            if (key==list[i]){
                return true;
            }
        }
        return false;
    }

    var task_id = mykeys.map(function(key, index) {
      if (key in self.state.data) {
        return buildFormGroup(key, self.state.data[key], 1);
      }
    });

    var tasks = remainder.map(function(key, index) {
      if (key in self.state.data) {
        return buildFormGroup(key, self.state.data[key], 1);
      }
    });

    var request_details = request_info.map(function(key, index) {
      if (key in self.state.data) {
          return buildFormGroup(key, self.state.data[key], 1);
      }
    });

    var extra = Object.keys(this.state.data).map(function(key, index) {
      if (!keyIn(key, mykeys.concat(remainder, request_info))) {
        return buildFormGroup(key, self.state.data[key], 1);
      }
    });

    var request = <div></div>;

    if (self.state.data["request_id"] !== undefined) {
      request = <TaskRequestItem requestid={self.state.data.request_id} />
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
        <h3>Task Host Details</h3>
        <table className="table table-striped">
            <tbody>
            {request_details}
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
        return buildFormGroup(key, self.state.data[key], 0);
      }
    });

    
    var ds = Object.keys(this.state.data).map(function(key, index) {
      if (keyIn(key, datakeys)) {
        return buildFormGroup(key, self.state.data[key], 0);
      }
    });
    
    var mc = Object.keys(this.state.data).map(function(key, index) {
      if (keyIn(key, memcachekeys)) {
        return buildFormGroup(key, self.state.data[key], 0);
      }
    });

    var extra = Object.keys(this.state.data).map(function(key, index) {
      if (!keyIn(key, memcachekeys) && !keyIn(key, datakeys) && !keyIn(key, originalkeys)) {
        return buildFormGroup(key, self.state.data[key], 0);
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

var buildFormGroup = function(prop, value, check) {

    function toTitleCase(str)
    {
        return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() 
            + txt.substr(1).toLowerCase();});
    }

    function get_480_120_60_10_5_1(value) {
        if (value > 480){
            return darkest;
        }
        if (value > 120){
            return darker;
        }
        if (value > 60){
            return dark;
        }
        if (value > 10){
            return light;
        }
        if (value > 5){
            return lighter;
        }
        if (value > 1){
            return lightest;
        }
        return " ";
    }

    function get_20_10_5(value) {
        if (value > 20){
            return darkest;
        }
        if (value > 10){
            return dark;
        }
        if (value > 5){
            return light;
        }
        return " ";
    }

    function get_120_60_30_10_1(value) {
        if (value > 120){
            return darkest;
        }
        if (value > 60){
            return dark;
        }
        if (value > 30){
            return light;
        }
        if (value > 10){
            return light;
        }
        if (value > 1){
            return light;
        }
        return " ";
    }

    var mystring = prop;
    mystring = mystring.replace(/_/g, ' ');
    mystring = toTitleCase(mystring);

    newValue = value.replace(/\D/g,'');

    var color = " ";

    if (prop == "ran" || prop == "task_eta" || prop == "end") {
      value = new Date(value*1000).toString();
    }

    switch (prop) {
        case "execution_count":
            if (value > 0){
                color = darkest;
            }
            break;
        case "retry_count":
            if (value > 0){
                color = darkest;
            }
            break;
        case "status_code":
            if (newValue != 200){
                color = darkest;
            }
            break;
        case "gae_latency_seconds":
            color = get_120_60_30_10_1(value);
            break;
        case "run_time":
            color = get_480_120_60_10_5_1(value);
            break;
        case "cpu_usage":
            if (newValue > 100){
                color = darkest;
            }
            break;
        case "memory_usage":
            if (newValue > 200){
                color = darkest;
            }
            break;
        case "rpc_total_count":
            if (newValue > 100){
                color = darkest;
            }
            break;
        case "taskqueue_BulkAdd_duration":
            if (newValue > 5){
                color = darkest;
            }
            break;
        case "exec_time":
            color = get_480_120_60_10_5_1(newValue);
            break;
        case "datastore_v3_Get_duration":
            color = get_20_10_5(newValue);
            break;
        case "datastore_v3_Put_duration":
            color = get_20_10_5(newValue);
            break;
        case "memcache_Delete_duration":
            color = get_20_10_5(newValue);
            break;
        case "memcache_Get_offset":
            color = get_20_10_5(newValue);
            break;
        case "memcache_Set_duration":
            color = get_20_10_5(newValue);
            break;
        case "urlfetch_Fetch_duration":
            color = get_20_10_5(newValue);
            break;
        default:
            color = " ";
    }

    if (check == 0) {
      return (
          <tr className={color}>
              <td title={prop} className="col-xs-4"><label>{mystring}:&nbsp;</label></td>
              <td className="col-xs-2"><span>{value}</span></td>
          </tr>
      )
    }
    return (
        <tr className={color}>
            <td title={prop} className="col-xs-2"><label>{mystring}:&nbsp;</label></td>
            <td className="col-xs-4"><span>{value}</span></td>
        </tr>
    )
}

module.exports = TaskDetail;
