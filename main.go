package main

import (
	"fmt"
	"os"

	"github.com/lyddonb/trajectory/trajectory"
)

func main() {
	// Redis test
	redisClient, err := trajectory.Connect()
	if err != err {
		fmt.Fprintf(os.Stderr, "Could not connection to redis: %s", err.Error())
	}

	trajectory.Listen(redisClient)
}
