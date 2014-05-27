package main

import (
	"net/http"

	"github.com/lyddonb/trajectory/db"
	"github.com/lyddonb/trajectory/rest"
)

const (
	TASK_PREFIX = "/tasks/"
)

func setupTasks(pool db.DBPool) {
	router := rest.SetupTaskRouter(pool, TASK_PREFIX)

	http.Handle(TASK_PREFIX, router)
}

func main() {
	// Stand up redis pool.
	pool := db.StartDB("127.0.0.1:6379", "")

	setupTasks(pool)

	http.ListenAndServe(":3000", nil)

	//listener := pipe.MakeConnection("tcp", ":1200")
	//taskPipeline := pipe.NewTaskPipeline(pool)
	//pipe.Listen(listener, taskPipeline)
}
