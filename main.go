package main

import (
	"github.com/lyddonb/trajectory/db"
	"github.com/lyddonb/trajectory/pipe"
)

func main() {
	// Stand up redis pool.
	pool := db.StartDB("127.0.0.1:6379", "")

	listener := pipe.MakeConnection("tcp", ":1200")
	taskPipeline := pipe.NewTaskPipeline(pool)
	pipe.Listen(listener, taskPipeline)
}
