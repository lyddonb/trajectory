package rest

import (
	"encoding/json"
	"net/http"

	"github.com/lyddonb/trajectory/api"
	"github.com/lyddonb/trajectory/db"
)

type TaskServices struct {
	api *api.TaskAPI
}

func NewTaskServices(taskAPI *api.TaskAPI) *TaskServices {
	return &TaskServices{taskAPI}
}

func (s *TaskServices) addTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var taskJson map[string]*json.RawMessage

	err := decoder.Decode(&taskJson)

	if err != nil {
		SendJsonErrorResponse(w, &taskJson, err.Error())
		return
	}

	task := api.ConvertTask(taskJson)

	_, ok := task[db.REQUEST_ADDRESS]

	if !ok {
		task[db.REQUEST_ADDRESS] = r.Host
	}

	timestamp, e := s.api.SaveTask(task)

	if e != nil {
		SendJsonErrorResponse(w, &taskJson, e.Error())
		return
	}

	SendJsonResponse(w, timestamp, nil)
}

func (s *TaskServices) getAllAddresses(w http.ResponseWriter, r *http.Request) {
	addresses, err := s.api.ListAddresses()

	SendJsonResponse(w, addresses, err)
}

func (s *TaskServices) getRequestsForAddress(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get(":address")

	if address == "" {
		SendJsonErrorResponse(w, nil, "No address passed in.")
		return
	}

	requests, err := s.api.ListRequests(address)

	SendJsonResponse(w, requests, err)
}

func (s *TaskServices) getTaskKeysForRequests(w http.ResponseWriter, r *http.Request) {
	requestid := r.URL.Query().Get(":requestid")

	if requestid == "" {
		SendJsonErrorResponse(w, nil, "No request id passed in.")
		return
	}

	taskKeys, err := s.api.ListRequestTaskKeys(requestid)

	SendJsonResponse(w, taskKeys, err)
}

func (s *TaskServices) getTaskKeysForTask(w http.ResponseWriter, r *http.Request) {
	taskid := r.URL.Query().Get(":taskid")

	if taskid == "" {
		SendJsonErrorResponse(w, nil, "No task id passed in.")
		return
	}

	taskKeys, err := s.api.ListTaskKeys(taskid)

	SendJsonResponse(w, taskKeys, err)
}

func (s *TaskServices) getTaskGraphForRequest(w http.ResponseWriter, r *http.Request) {
	requestid := r.URL.Query().Get(":requestid")

	if requestid == "" {
		SendJsonErrorResponse(w, nil, "No request id passed in.")
		return
	}

	graph, err := s.api.GetRequestTaskGraph(requestid)

	SendJsonResponse(w, graph, err)
}

func (s *TaskServices) getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Send it"))
}

func (s *TaskServices) getTaskByKey(w http.ResponseWriter, r *http.Request) {
	taskKey := r.URL.Query().Get(":taskKey")

	if taskKey == "" {
		SendJsonErrorResponse(w, nil, "No task key passed in.")
		return
	}

	task, err := s.api.GetTaskForKey(taskKey)

	SendJsonResponse(w, task, err)
}
