package api

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/lyddonb/trajectory/db"
)

func TestListTaskRequests(t *testing.T) {
	pool := db.StartDB("127.0.0.1:6379", "")
	timestamp := strconv.FormatInt(time.Now().UTC().Unix(), 10)

	conn := pool.Get()
	conn.Send("ZADD", PARENT_REQUESTS, timestamp, "parentreqeustid")
	conn.Close()

	result := ListTaskRequests(pool)

	fmt.Println(string(result))
}
