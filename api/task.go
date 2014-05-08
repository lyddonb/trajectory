package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/lyddonb/trajectory/db"
)

const (
	PARENT_REQUESTS = "ParentRequests"
)

// TODO: Have the take requests take a "machine" id to filter by.
// TODO: Create a list of "machines" that we are tracking from.

// Returns a byte array of a jsonified list of strings (request ids).
func ListTaskRequests(pool db.DBPool) []byte {
	conn := pool.Get()
	defer conn.Close()

	set := make(map[string]int)

	r, _ := redis.Strings(conn.Do("ZREVRANGE", PARENT_REQUESTS, 0, -1, "WITHSCORES"))
	fmt.Printf("Response from redis %s.\n", r)

	var lastItem string

	for count, item := range r {
		if item == "" {
			lastItem = ""
			continue
		}

		// TODO: Maybe check if the parent request id already exists and make it
		// a list of timestamps.
		if (count % 2) == 0 {
			lastItem = string(item)
			set[lastItem] = 0
		} else if lastItem != "" {
			val, err := strconv.Atoi(item)
			if err != nil {
				set[lastItem] = 0
			} else {
				set[lastItem] = val
			}
		}
	}

	// TODO: filter empty and duplicate strings.

	b, _ := json.Marshal(set)

	return b
}
