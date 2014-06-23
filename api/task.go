package api

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/lyddonb/trajectory/db"
)

// TODO: Have the take requests take a "machine" id to filter by.
// TODO: Create a list of "machines" that we are tracking from.

type TaskAPI struct {
	dal db.TaskDAL
}

func NewTaskAPI(pool db.DBPool) *TaskAPI {
	return &TaskAPI{db.NewTaskDataAccess(pool)}
}

func convertWeightedListToSet(taskScores []string) map[string]int {
	var lastItem string
	set := make(map[string]int)

	for count, item := range taskScores {
		if item == "" {
			lastItem = ""
			continue
		}

		if (count % 2) == 0 {
			lastItem = string(item)
			set[lastItem] = 0
		} else if lastItem != "" {
			val, err := strconv.Atoi(item)
			// QUESTION: What to do with the error?
			if err != nil {
				set[lastItem] = 0
			} else {
				set[lastItem] = val
			}
		}
	}

	return set
}

func (a *TaskAPI) SaveTask(task db.Task) (string, error) {
	splitKey := strings.Split(task["id"], ":")

	if len(splitKey) == 2 {
		task[db.PARENT_TASK_ID] = ""
		task[db.PARENT_REQUEST_ID] = splitKey[0]
		task[db.TASK_ID] = splitKey[1]
	} else {
		task[db.PARENT_TASK_ID] = splitKey[1]
		task[db.PARENT_REQUEST_ID] = splitKey[0]
		task[db.TASK_ID] = splitKey[2]
	}

	return a.dal.SaveTask(task)
}

func ConvertTask(taskJson map[string]*json.RawMessage) db.Task {
	taskMap := make(db.Task)

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

// Returns a byte array of a jsonified list of strings (request ids).
func (a *TaskAPI) ListRequests(address string) (map[string]int, error) {
	taskScores, error := a.dal.GetRequests(address)

	if error != nil {
		return nil, error
	}

	return convertWeightedListToSet(taskScores), nil
}

func (a *TaskAPI) ListAddresses() (map[string]int, error) {
	addressScores, error := a.dal.GetAddresses()

	if error != nil {
		return nil, error
	}

	return convertWeightedListToSet(addressScores), nil
}

func (a *TaskAPI) ListRequestTaskKeys(requestId string) (map[string]int, error) {
	taskKeyScores, err := a.dal.GetRequestTaskKeys(requestId)

	if err != nil {
		return nil, err
	}

	return convertWeightedListToSet(taskKeyScores), nil
}

//func ListRequestTasks(taskData db.DataAccess) ([]db.Task, error) {
//}

func (a *TaskAPI) GetTaskForKey(taskKey string) (db.Task, error) {
	return a.dal.GetTaskForKey(taskKey)
}
