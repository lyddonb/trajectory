package db

import (
	"fmt"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
)

type TaskDAL interface {
	GetRequests(machine string) ([]string, error)
	GetRequestTaskKeys(requestId string) ([]string, error)
	GetAddresses() ([]string, error)
	SaveTask(task Task) (string, error)
	GetTaskForKey(taskKey string) (Task, error)
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

func (c *TaskDataAccess) getRangeResults(key string) ([]string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	return redis.Strings(
		conn.Do("ZREVRANGE", key, 0, -1, "WITHSCORES"))
}

func (c *TaskDataAccess) GetRequests(address string) ([]string, error) {
	return c.getRangeResults(getRequestKey(address))
}

func (c *TaskDataAccess) GetAddresses() ([]string, error) {
	address, err := c.getRangeResults(ADDRESS)

	return address, err
}

func (c *TaskDataAccess) GetRequestTaskKeys(requestId string) ([]string, error) {
	return c.getRangeResults(requestId)
}

func getRequestKey(address string) string {
	return fmt.Sprintf("%s:%s", PARENT_REQUESTS, address)
}

func (c *TaskDataAccess) GetTaskForKey(taskKey string) (Task, error) {
	conn := c.pool.Get()
	defer conn.Close()

	v, err := redis.Values(conn.Do("HGETALL", taskKey))

	if err != nil {
		return nil, err
	}

	if v == nil {
		return nil, nil
	}

	return ScanMap(v)
}

func (c *TaskDataAccess) SaveTask(task Task) (string, error) {
	conn := c.pool.Get()
	defer conn.Close()

	taskKey := task.Key()

	if taskKey == "::" {
		fmt.Println("Not a valid task. %s", task)
		// TODO: Make this return an error.
		return "", nil
	}

	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	parentRequestId := task[PARENT_REQUEST_ID]

	err := AddTask(conn, taskKey, task)

	if err != nil {
		return "", err
	}

	err = AddTaskToParentRequest(conn, parentRequestId, timestamp, taskKey)

	if err != nil {
		return "", err
	}

	err = AddParentRequest(conn, task[REQUEST_ADDRESS], timestamp, parentRequestId)

	if err != nil {
		return "", err
	}

	err = AddAddress(conn, task[REQUEST_ADDRESS], timestamp)

	if err != nil {
		return "", err
	}

	err = conn.Flush()

	if err != nil {
		return "", err
	}

	return timestamp, nil
}

func AddTask(conn redis.Conn, taskKey string, task Task) error {
	return conn.Send("HMSET", redis.Args{}.Add(taskKey).AddFlat(task)...)
}

func AddTaskToParentRequest(
	conn redis.Conn, parentRequestId, timestamp, taskKey string) error {
	return conn.Send("ZADD", parentRequestId, timestamp, taskKey)
}

func AddParentRequest(conn redis.Conn, address, timestamp, parentRequestId string) error {
	return conn.Send("ZADD", getRequestKey(address), timestamp, parentRequestId)
}

func AddAddress(conn redis.Conn, address, timestamp string) error {
	return conn.Send("ZADD", ADDRESS, timestamp, address)
}
