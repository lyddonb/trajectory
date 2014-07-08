package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lyddonb/trajectory/db"
	"github.com/lyddonb/trajectory/pipe"
	"github.com/lyddonb/trajectory/rest"
)

const (
	TASK_PREFIX = "/api/tasks/"
	STAT_PREFIX = "/api/stats/"
)

func setupTasks(pool db.DBPool, writeToFile bool) {
	router := rest.SetupTaskRouter(pool, TASK_PREFIX, writeToFile)

	http.Handle(TASK_PREFIX, router)
}

func setupStats(pool db.DBPool, writeToFile bool) {
	router := rest.SetupStatRouter(pool, STAT_PREFIX, writeToFile)

	http.Handle(STAT_PREFIX, router)
}

func setupWeb() {
	router := mux.NewRouter()
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	http.Handle("/", router)
}

func main() {
	// Stand up redis pool.
	pool := db.StartDB("127.0.0.1:6379", "")

	go func() {
		listener := pipe.MakeConnection("tcp", ":1300")
		taskPipeline := pipe.NewTaskPipeline(pool)
		pipe.Listen(listener, taskPipeline)
	}()

	writeToFile := false

	setupStats(pool, writeToFile)
	setupTasks(pool, writeToFile)
	setupWeb()

	http.ListenAndServe(":3000", nil)
}
