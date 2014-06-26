package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
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
	//var item map[string]string
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
	params := mux.Vars(r)
	address := params["address"]

	if address == "" {
		SendJsonErrorResponse(w, nil, "No address passed in.")
	}

	requests, err := s.api.ListRequests(address)

	SendJsonResponse(w, requests, err)
}

func (s *TaskServices) getTaskKeysForRequests(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	requestid := params["requestid"]

	if requestid == "" {
		SendJsonErrorResponse(w, nil, "No request id passed in.")
	}

	taskKeys, err := s.api.ListRequestTaskKeys(requestid)

	SendJsonResponse(w, taskKeys, err)
}

func (s *TaskServices) getTaskGraphForRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	requestid := params["requestid"]

	if requestid == "" {
		SendJsonErrorResponse(w, nil, "No request id passed in.")
	}

	graph, err := s.api.GetRequestTaskGraph(requestid)

	SendJsonResponse(w, graph, err)
}

func (s *TaskServices) getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Send it"))
}

func (s *TaskServices) getTaskByKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskKey := params["taskKey"]

	if taskKey == "" {
		SendJsonErrorResponse(w, nil, "No task key passed in.")
	}

	task, err := s.api.GetTaskForKey(taskKey)

	SendJsonResponse(w, task, err)
}
