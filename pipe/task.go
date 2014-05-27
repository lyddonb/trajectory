package pipe

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/lyddonb/trajectory/api"
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

	err := json.Unmarshal(message, &taskJson)

	if err != nil {
		fmt.Println("Error parsing task json %s", err)
		return make(db.Task)
	}

	return api.ConvertTask(taskJson)
}
