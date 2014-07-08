package api

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"

	"github.com/lyddonb/trajectory/db"
)

// TODO: Have the take requests take a "machine" id to filter by.
// TODO: Create a list of "machines" that we are tracking from.

type TaskAPI struct {
	dal db.TaskDAL
}

// A data structure to hold a key/value pair.
type WeightedScore struct {
	Key   string
	Value int
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type WeightedScoreList []WeightedScore

func (p WeightedScoreList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p WeightedScoreList) Len() int           { return len(p) }
func (p WeightedScoreList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a WeightedScoreList, then sort and return it.
func sortMapByValue(m map[string]int) WeightedScoreList {
	p := make(WeightedScoreList, len(m))
	i := 0
	for k, v := range m {
		p[i] = WeightedScore{k, v}
		i++
	}
	sort.Sort(sort.Reverse(p))
	return p
}

func NewTaskAPI(pool db.DBPool) *TaskAPI {
	return &TaskAPI{db.NewTaskDataAccess(pool)}
}

func mergeWeightedList(taskScores []string) WeightedScoreList {
	mapScores := convertWeightedListToSet(taskScores)

	return sortMapByValue(mapScores)
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
func (a *TaskAPI) ListRequests(address string) (WeightedScoreList, error) {
	taskScores, error := a.dal.GetRequests(address)

	if error != nil {
		return nil, error
	}

	return mergeWeightedList(taskScores), nil
}

func (a *TaskAPI) ListAddresses() (WeightedScoreList, error) {
	addressScores, error := a.dal.GetAddresses()

	if error != nil {
		return nil, error
	}

	return mergeWeightedList(addressScores), nil
}

func (a *TaskAPI) ListRequestTaskKeys(requestId string) (WeightedScoreList, error) {
	taskKeyScores, err := a.dal.GetRequestTaskKeys(requestId)

	if err != nil {
		return nil, err
	}

	return mergeWeightedList(taskKeyScores), nil
}

func (a *TaskAPI) ListRequestTaskKeysAsSet(requestId string) (map[string]int, error) {
	taskKeyScores, err := a.dal.GetRequestTaskKeys(requestId)

	if err != nil {
		return nil, err
	}

	return convertWeightedListToSet(taskKeyScores), nil
}

func (a *TaskAPI) ListTaskKeys(taskId string) (WeightedScoreList, error) {
	taskKeyScores, err := a.dal.GetTaskKeys(taskId)

	if err != nil {
		return nil, err
	}

	return mergeWeightedList(taskKeyScores), nil
}

//func ListRequestTasks(taskData db.DataAccess) ([]db.Task, error) {
//}

func (a *TaskAPI) GetTaskForKey(taskKey string) (db.Task, error) {
	return a.dal.GetTaskForKey(taskKey)
}
