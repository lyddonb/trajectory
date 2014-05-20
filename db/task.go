package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

const (
	PARENT_TASK_ID    = "parent_task_id"
	TASK_ID           = "task_id"
	REQUEST_ID        = "request_id"
	PARENT_REQUESTS   = "ParentRequests"
	PARENT_REQUEST_ID = "parent_request_id"
)

type DataAccess interface {
	GetRequests(machine string) ([]string, error)
}

type TaskDataAccess struct {
	pool DBPool
}

type Task map[string]string

func (t Task) Key() string {
	return fmt.Sprintf("%s:%s:%s", t[PARENT_TASK_ID], t[TASK_ID], t[REQUEST_ID])
}

func NewTaskDataAccess(pool DBPool) *TaskDataAccess {
	return &TaskDataAccess{pool}
}

func (c *TaskDataAccess) GetRequests(machine string) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Strings(
		conn.Do("ZREVRANGE", PARENT_REQUESTS, 0, -1, "WITHSCORES"))
}

func (c *TaskDataAccess) SaveTask(task Task) string {
	conn := c.pool.Get()
	defer conn.Close()

	taskKey := task.Key()

	if taskKey == "::" {
		fmt.Println("Not a valid task. %s", task)
		return ""
	}

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	parentRequestId := task[PARENT_REQUEST_ID]

	AddTask(conn, taskKey, task)
	AddTaskToParentRequest(conn, parentRequestId, timestamp, taskKey)
	AddParentRequest(conn, timestamp, parentRequestId)

	conn.Flush()

	return timestamp
}

func AddTask(conn redis.Conn, taskKey string, task Task) {
	conn.Send("HMSET", redis.Args{taskKey}.AddFlat(task)...)
}

func AddTaskToParentRequest(
	conn redis.Conn, parentRequestId, timestamp, taskKey string) {
	conn.Send("ZADD", parentRequestId, timestamp, taskKey)
}

func AddParentRequest(conn redis.Conn, timestamp, parentRequestId string) {
	conn.Send("ZADD", PARENT_REQUESTS, timestamp, parentRequestId)
}
