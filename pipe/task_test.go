package pipe

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseTaskWithString(t *testing.T) {
	jsonObj := map[string]string{
		"propa": "a",
	}

	b, _ := json.Marshal(jsonObj)

	result := ParseTask(b)

	assert.Equal(t, result["propa"], "a")
}

func TestParseTaskWithInt(t *testing.T) {
	jsonObj := map[string]int{
		"propa": 3,
	}

	b, _ := json.Marshal(jsonObj)

	result := ParseTask(b)

	assert.Equal(t, result["propa"], "3")
}

func TestParseTaskWithInt64(t *testing.T) {
	jsonObj := map[string]int64{
		"propa": 9223372036854775807,
	}

	b, _ := json.Marshal(jsonObj)

	result := ParseTask(b)

	assert.Equal(t, result["propa"], "9223372036854775807")
}

func TestParseTaskWithFloat(t *testing.T) {
	jsonObj := map[string]float64{
		"propa": 203685477580.2343,
	}

	b, _ := json.Marshal(jsonObj)

	result := ParseTask(b)

	assert.Equal(t, result["propa"], "203685477580.2343")
}

func TestParseTaskWithBool(t *testing.T) {
	jsonObj := map[string]bool{
		"propa": true,
	}

	b, _ := json.Marshal(jsonObj)

	result := ParseTask(b)

	assert.Equal(t, result["propa"], "true")
}

type MockRedisPool struct {
	mock.Mock
	DBPool
}

func (p *MockRedisPool) Get() redis.Conn {
	args := p.Mock.Called()
	return args.Get(0).(redis.Conn)
}

type MockRedisConnection struct {
	mock.Mock
	redis.Conn
}

func (c *MockRedisConnection) Close() error {
	args := c.Mock.Called()
	return args.Error(0)
}

func verifyArgs(arg, value string) {
	if arg != value {
		fmt.Println("Failed %s did not equal %s", arg, value)
	}
}

func (c *MockRedisConnection) Send(commandName string, args ...interface{}) error {
	funcArgs := c.Mock.Called()

	if commandName == "HMSET" {
		verifyArgs(args[0].(string), "parenttaskid:taskid:request_id")
		verifyArgs(args[1].(string), "parent_task_id")
		verifyArgs(args[2].(string), "parenttaskid")
		verifyArgs(args[3].(string), "parent_request_id")
		verifyArgs(args[4].(string), "parentrequestid")
		verifyArgs(args[5].(string), "task_id")
		verifyArgs(args[6].(string), "taskid")
		verifyArgs(args[7].(string), "request_id")
		verifyArgs(args[8].(string), "request_id")
	} else if commandName == "ZADD" && args[0] == PARENT_REQUESTS {
		verifyArgs(args[2].(string), "parentrequestid")
	} else {
		verifyArgs(args[0].(string), "parentrequestid")
		verifyArgs(args[2].(string), "parenttaskid:taskid:request_id")
	}

	return funcArgs.Error(0)
}

func (c *MockRedisConnection) Flush() error {
	args := c.Mock.Called()
	return args.Error(0)
}

func TestWriteTask(t *testing.T) {
	redisPool := new(MockRedisPool)
	redisConn := new(MockRedisConnection)

	task := Task{
		PARENT_TASK_ID:    "parenttaskid",
		PARENT_REQUEST_ID: "parentrequestid",
		TASK_ID:           "taskid",
		REQUEST_ID:        "request_id",
	}

	redisPool.On("Get").Return(redisConn)

	redisConn.On("Close").Return(nil)

	redisConn.On("Send").Return(nil)
	redisConn.On("Send").Return(nil)
	redisConn.On("Send").Return(nil)

	redisConn.On("Flush").Return(nil)

	WriteTask(task, redisPool)

	redisConn.Mock.AssertExpectations(t)
	redisPool.Mock.AssertExpectations(t)

	redisConn.AssertCalled(t, "Send")
}
