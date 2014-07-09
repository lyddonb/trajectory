package rest

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/lyddonb/trajectory/api"
	"github.com/lyddonb/trajectory/db"
)

func SetupTasks(pool db.DBPool) *TaskServices {
	return NewTaskServices(api.NewTaskAPI(pool))
}

func SetupStats(pool db.DBPool) *StatServices {
	return NewStatServices(api.NewStatAPI(pool))
}

type RestMiddleware func(http.HandlerFunc, ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc

func SetupTaskRouter(pool db.DBPool, prefix string, writeToFile bool, middleware RestMiddleware) *pat.PatternServeMux {
	router := pat.New()

	taskServices := SetupTasks(pool)

	router.Get(prefix+"addresses", middleware(taskServices.getAllAddresses))
	router.Get(prefix+"addresses/:address/requests",
		middleware(taskServices.getRequestsForAddress))
	router.Get(prefix+"addresses/:address/requests/:requestid/tasks",
		middleware(taskServices.getTaskKeysForRequests))
	router.Get(prefix+"addresses/:address/requests/:requestid/taskgraph",
		middleware(taskServices.getTaskGraphForRequest))
	router.Get(prefix+"task/:taskKey", middleware(taskServices.getTaskByKey))
	router.Get(prefix+"tasks/:taskid",
		middleware(taskServices.getTaskKeysForTask))

	if writeToFile {
		//router.HandleFunc(prefix, middleware(taskServices.addTaskToFile)).Methods("POST")
		router.Post(prefix, middleware(taskServices.addTaskToFile))
	} else {
		router.Post(prefix, http.HandlerFunc(taskServices.addTask))
	}
	router.Post(prefix, http.HandlerFunc(taskServices.addTask))
	router.Get(prefix, middleware(taskServices.getAllTasks))

	return router
}

func SetupStatRouter(pool db.DBPool, prefix string, writeToFile bool, middleware RestMiddleware) *pat.PatternServeMux {
	router := pat.New()

	statServices := SetupStats(pool)

	router.Get(prefix+":requestId", middleware(statServices.getStatByRequestId))
	if writeToFile {
		router.Post(prefix, http.HandlerFunc(statServices.addStatToFile))
	} else {
		router.Post(prefix, http.HandlerFunc(statServices.addStat))
	}
	router.Get(prefix, middleware(statServices.getAllStats))

	return router
}
