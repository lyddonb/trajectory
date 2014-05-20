package api

import (
	"strconv"

	"github.com/lyddonb/trajectory/db"
)

// TODO: Have the take requests take a "machine" id to filter by.
// TODO: Create a list of "machines" that we are tracking from.

type TaskAPI struct {
	dal db.DataAccess
}

func NewTaskAPI(dal db.DataAccess) *TaskAPI {
	return &TaskAPI{dal}
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
	taskKeyScores, error := a.dal.GetRequestTaskKeys(requestId)

	if error != nil {
		return nil, error
	}

	return convertWeightedListToSet(taskKeyScores), nil
}

//func ListRequestTasks(taskData db.DataAccess) ([]db.Task, error) {
//}

//func GetTaskForKey(taskData db.DataAccess) (db.Task, error) {
//}
