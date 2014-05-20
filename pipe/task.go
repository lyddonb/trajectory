package pipe

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/lyddonb/trajectory/db"
)

type TaskPipeline struct {
	isOpen bool
	dal    *db.TaskDataAccess
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

func (tp *TaskPipeline) Parse(message []byte, remoteAddr string) {
	task := ParseTask(message)

	task[db.REQUEST_ADDRESS] = remoteAddr

	tp.dal.SaveTask(task)
}

func NewTaskPipeline(pool db.DBPool) *TaskPipeline {
	dal := db.NewTaskDataAccess(pool)

	return &TaskPipeline{
		isOpen: true,
		dal:    dal,
	}
}

func ParseTask(message []byte) db.Task {
	var taskJson map[string]*json.RawMessage
	taskMap := make(db.Task)

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
