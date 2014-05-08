package pipe

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	PARENT_TASK_ID    = "parent_task_id"
	TASK_ID           = "task_id"
	REQUEST_ID        = "request_id"
	PARENT_REQUEST_ID = "parent_request_id"
	PARENT_REQUESTS   = "ParentRequests"
)

type Task map[string]string

func (t Task) Key() string {
	return fmt.Sprintf("%s:%s:%s", t[PARENT_TASK_ID], t[TASK_ID], t[REQUEST_ID])
}

type TaskPipeline struct {
	isOpen bool
}

func (tp *TaskPipeline) Handler(conn net.Conn) {
	handleClient(conn, tp)
}

func (tp *TaskPipeline) Error(err error) {
}

func (tp *TaskPipeline) Open() bool {
	// TODO: Handle a global variable or channel to close the connection.
	return tp.isOpen
}

func (tp *TaskPipeline) Parse(message []byte) {
	task := ParseTask(message)
	WriteTask(task, pool)
}

func NewTaskPipeline() *TaskPipeline {
	return &TaskPipeline{isOpen: true}
}

func ParseTask(message []byte) Task {
	var taskJson map[string]*json.RawMessage
	taskMap := make(Task)

	err := json.Unmarshal(message, &taskJson)

	if err != nil {
		fmt.Println("Error parsing task json %s", err)
		return taskMap
	}

	for key, value := range taskJson {
		var stringValue string
		var intValue int
		var floatValue float64
		var boolValue bool

		if json.Unmarshal(*value, &stringValue) == nil {
			taskMap[key] = stringValue
		} else if json.Unmarshal(*value, &intValue) == nil {
			taskMap[key] = strconv.Itoa(intValue)
		} else if json.Unmarshal(*value, &floatValue) == nil {
			taskMap[key] = strconv.FormatFloat(floatValue, 'f', -1, 64)
		} else if json.Unmarshal(*value, &boolValue) == nil {
			taskMap[key] = strconv.FormatBool(boolValue)
		}
	}

	return taskMap
}

func WriteTask(task Task, redisPool DBPool) string {
	fmt.Println("Write Task.")
	conn := redisPool.Get()
	defer conn.Close()

	taskKey := task.Key()

	if taskKey == "::" {
		fmt.Println("Not a valid task. %s", task)
		return ""
	}

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	parentRequestId := task[PARENT_REQUEST_ID]

	conn.Send("HMSET", redis.Args{taskKey}.AddFlat(task)...)
	conn.Send("ZADD", parentRequestId, timestamp, taskKey)
	conn.Send("ZADD", PARENT_REQUESTS, timestamp, parentRequestId)

	conn.Flush()

	r, _ := conn.Do("ZREVRANGE", PARENT_REQUESTS, 0, 150)
	fmt.Printf("Response from redis %s.\n", r)

	return timestamp
}
