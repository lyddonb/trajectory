package rest

import (
	"github.com/gorilla/mux"
	"github.com/lyddonb/trajectory/api"
	"github.com/lyddonb/trajectory/db"
)

func SetupTasks(pool db.DBPool) *TaskServices {
	return NewTaskServices(api.NewTaskAPI(pool))
}

func SetupStats(pool db.DBPool) *StatServices {
	return NewStatServices(api.NewStatAPI(pool))
}

func SetupTaskRouter(pool db.DBPool, prefix string, writeToFile bool) *mux.Router {
	router := mux.NewRouter()

	taskServices := SetupTasks(pool)

	if writeToFile {
		router.HandleFunc(prefix, taskServices.addTaskToFile).Methods("POST")
	} else {
		router.HandleFunc(prefix, taskServices.addTask).Methods("POST")
	}
	router.HandleFunc(prefix, taskServices.addTask).Methods("POST")
	router.HandleFunc(prefix, taskServices.getAllTasks).Methods("GET")
	router.HandleFunc(prefix+"addresses",
		taskServices.getAllAddresses).Methods("GET")
	router.HandleFunc(prefix+"addresses/{address}/requests",
		taskServices.getRequestsForAddress).Methods("GET")
	router.HandleFunc(prefix+"addresses/{address}/requests/{requestid}/tasks",
		taskServices.getTaskKeysForRequests).Methods("GET")
	router.HandleFunc(prefix+"addresses/{address}/requests/{requestid}/taskgraph",
		taskServices.getTaskGraphForRequest).Methods("GET")
	router.HandleFunc(prefix+"task/{taskKey}",
		taskServices.getTaskByKey).Methods("GET")
	router.HandleFunc(prefix+"tasks/{taskid}",
		taskServices.getTaskKeysForTask).Methods("GET")

	return router
}

func SetupStatRouter(pool db.DBPool, prefix string, writeToFile bool) *mux.Router {
	router := mux.NewRouter()

	statServices := SetupStats(pool)

	if writeToFile {
		router.HandleFunc(prefix, statServices.addStatToFile).Methods("POST")
	} else {
		router.HandleFunc(prefix, statServices.addStat).Methods("POST")
	}
	router.HandleFunc(prefix, statServices.getAllStats).Methods("GET")
	router.HandleFunc(prefix+"{requestId}", statServices.getStatByRequestId).Methods("GET")

	return router
}
