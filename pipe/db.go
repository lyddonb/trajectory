package pipe

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

type DBPool interface {
	Get() redis.Conn
}

// Taken from: http://godoc.org/github.com/garyburd/redigo/redis#Pool
func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func StartDB(server, password string) {
	pool = newPool(server, password)
}
