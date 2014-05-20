package api

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/lyddonb/trajectory/db"
	"github.com/stretchr/testify/assert"
)

func (p *MockTaskDataAccess) GetTaskForKey(taskKey string) (db.Task, error) {
	args := p.Mock.Called()

	var task db.Task

	for _, arg := range args {
		task = arg.(db.Task)

		if task.Key() == taskKey {
			return task, nil
		}
	}

	return args.Get(0).(db.Task), nil
}

func TestParentOnlyGraph(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	requestId := "requestId"

	taskKeyResult := []string{}

	dal.On("GetRequestTaskKeys").Return(taskKeyResult, nil)

	parentNode, _ := api.GetRequestTaskGraph(requestId)

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, parentNode.TaskId, requestId)
	assert.Equal(t, parentNode.Name, "Parent Request")
	assert.Equal(t, parentNode.ContextId, "")
	assert.Equal(t, len(parentNode.Children), 0)
	assert.Equal(t, len(parentNode.Keys), 0)
	assert.True(t, parentNode.IsParent)
}

func TestParentWithSingleNodeGraph(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	parentRequestId := "parentrequestid"
	requestId := "requestId"

	timestampStr := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	childTaskid := "taskid"
	childTaskKey := fmt.Sprintf("%s:%s:%s", parentRequestId, childTaskid,
		requestId)
	childUrl := "a.test.url"

	taskKeyResult := []string{
		childTaskKey,
		timestampStr,
	}

	task := makeTask(parentRequestId, parentRequestId, childTaskid, requestId, childUrl)

	dal.On("GetRequestTaskKeys").Return(taskKeyResult, nil)
	dal.On("GetTaskForKey").Return(task, nil)

	parentNode, _ := api.GetRequestTaskGraph(parentRequestId)

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, parentNode.TaskId, parentRequestId)
	assert.Equal(t, parentNode.Name, "Parent Request")
	assert.Equal(t, parentNode.ContextId, "")
	assert.Equal(t, len(parentNode.Children), 1)
	assert.Equal(t, len(parentNode.Keys), 0)
	assert.True(t, parentNode.IsParent)

	childNode, _ := parentNode.Children[childTaskKey]

	assert.Equal(t, childNode.TaskId, childTaskid)
	assert.Equal(t, childNode.Name, childUrl)
	assert.Equal(t, childNode.ContextId, "")
	assert.Equal(t, len(childNode.Children), 0)
	assert.Equal(t, len(childNode.Keys), 1)
	assert.Equal(t, childNode.Keys[0], childTaskKey)
	assert.False(t, childNode.IsParent)
}

func TestParentWithTwoChildNodesGraph(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	parentRequestId := "parentrequestid"
	requestId := "requestId"
	requestId2 := "requestId2"

	contextid := "contextid"
	childTaskid := "taskid"
	child2Taskid := "taskid2"

	childTaskKey := fmt.Sprintf("%s:%s:%s", parentRequestId,
		childTaskid+"|"+contextid, requestId)
	child2TaskKey := fmt.Sprintf("%s:%s:%s", parentRequestId,
		child2Taskid+"|"+contextid, requestId2)

	childUrl := "a.test.url"
	child2Url := "a.second.test.url"

	taskKeyResult := []string{
		childTaskKey,
		"1",
		child2TaskKey,
		"2",
	}

	task := makeTask(parentRequestId, parentRequestId,
		childTaskid+"|"+contextid, requestId, childUrl)
	task2 := makeTask(parentRequestId, parentRequestId,
		child2Taskid+"|"+contextid, requestId2, child2Url)

	dal.On("GetRequestTaskKeys").Return(taskKeyResult, nil)
	dal.On("GetTaskForKey").Return(task, task2).Times(2)

	parentNode, _ := api.GetRequestTaskGraph(parentRequestId)

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, parentNode.TaskId, parentRequestId)
	assert.Equal(t, parentNode.Name, "Parent Request")
	assert.Equal(t, parentNode.ContextId, "")
	assert.Equal(t, len(parentNode.Children), 2)
	assert.Equal(t, len(parentNode.Keys), 0)
	assert.True(t, parentNode.IsParent)

	childNode, _ := parentNode.Children[childTaskKey]

	assert.Equal(t, childNode.TaskId, childTaskid)
	assert.Equal(t, childNode.Name, childUrl)
	assert.Equal(t, childNode.ContextId, contextid)
	assert.Equal(t, len(childNode.Children), 0)
	assert.Equal(t, len(childNode.Keys), 1)
	assert.Equal(t, childNode.Keys[0], childTaskKey)
	assert.False(t, childNode.IsParent)

	childNode2, _ := parentNode.Children[child2TaskKey]
	assert.Equal(t, childNode2.TaskId, child2Taskid)
	assert.Equal(t, childNode2.Name, child2Url)
	assert.Equal(t, childNode2.ContextId, contextid)
	assert.Equal(t, len(childNode2.Children), 0)
	assert.Equal(t, len(childNode2.Keys), 1)
	assert.Equal(t, childNode2.Keys[0], child2TaskKey)
	assert.False(t, childNode2.IsParent)
}

func TestParentWithChildNodeWithChildNodeGraph(t *testing.T) {
	dal := new(MockTaskDataAccess)
	api := NewTaskAPI(dal)

	parentRequestId := "parentrequestid"
	requestId := "requestId"
	requestId2 := "requestId2"

	childTaskid := "taskid"
	child2Taskid := "taskid2"

	childTaskKey := fmt.Sprintf("%s:%s:%s", parentRequestId, childTaskid, requestId)

	child2TaskKey := fmt.Sprintf("%s:%s:%s", childTaskid, child2Taskid, requestId2)

	childUrl := "a.test.url"
	child2Url := "a.second.test.url"

	taskKeyResult := []string{
		childTaskKey,
		"1",
		child2TaskKey,
		"2",
	}

	task := makeTask(parentRequestId, parentRequestId,
		childTaskid, requestId, childUrl)
	task2 := makeTask(parentRequestId, childTaskid,
		child2Taskid, requestId2, child2Url)

	dal.On("GetRequestTaskKeys").Return(taskKeyResult, nil)
	dal.On("GetTaskForKey").Return(task, task2).Times(2)

	parentNode, _ := api.GetRequestTaskGraph(parentRequestId)

	dal.Mock.AssertExpectations(t)

	assert.Equal(t, parentNode.TaskId, parentRequestId)
	assert.Equal(t, parentNode.Name, "Parent Request")
	assert.Equal(t, parentNode.ContextId, "")
	assert.Equal(t, len(parentNode.Children), 1)
	assert.Equal(t, len(parentNode.Keys), 0)
	assert.True(t, parentNode.IsParent)

	childNode, _ := parentNode.Children[childTaskKey]

	assert.Equal(t, childNode.TaskId, childTaskid)
	assert.Equal(t, childNode.Name, childUrl)
	assert.Equal(t, childNode.ContextId, "")
	assert.Equal(t, len(childNode.Children), 1)
	assert.Equal(t, len(childNode.Keys), 1)
	assert.Equal(t, childNode.Keys[0], childTaskKey)
	assert.False(t, childNode.IsParent)

	childNode2, _ := childNode.Children[child2TaskKey]
	assert.Equal(t, childNode2.TaskId, child2Taskid)
	assert.Equal(t, childNode2.Name, child2Url)
	assert.Equal(t, childNode2.ContextId, "")
	assert.Equal(t, len(childNode2.Children), 0)
	assert.Equal(t, len(childNode2.Keys), 1)
	assert.Equal(t, childNode2.Keys[0], child2TaskKey)
	assert.False(t, childNode2.IsParent)
}

func makeTask(parentRequestId, parentTaskId, childTaskid, requestId, url string) db.Task {
	task := make(db.Task)
	task[db.PARENT_REQUEST_ID] = parentRequestId
	task[db.PARENT_TASK_ID] = parentTaskId
	task[db.TASK_ID] = childTaskid
	task[db.REQUEST_ID] = requestId
	task[db.URL] = url

	return task
}
