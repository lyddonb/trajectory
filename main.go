package main

import (
	"fmt"
	"os"

	"github.com/lyddonb/trajectory/trajectory"
)

func main() {
	// Web host
	go func() {
		trajectory.StartWeb()
	}()

	go func() {
		trajectory.StartRest()
	}()

	// Redis test
	redisClient, err := trajectory.Connect()
	if err != err {
		fmt.Fprintf(os.Stderr, "Could not connection to redis: %s", err.Error())
	}

	trajectory.Listen(redisClient)
}
