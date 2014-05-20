package api

import (
	"strconv"

	"github.com/lyddonb/trajectory/db"
)

// TODO: Have the take requests take a "machine" id to filter by.
// TODO: Create a list of "machines" that we are tracking from.

func convertTaskListToSet(taskScores []string) map[string]int {
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
func ListTaskRequests(taskData db.DataAccess, machine string) (map[string]int, error) {
	taskScores, error := taskData.GetRequests(machine)

	if error != nil {
		return nil, error
	}

	return convertTaskListToSet(taskScores), nil
}
