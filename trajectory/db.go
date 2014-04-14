package trajectory

import (
	"os"

	"github.com/garyburd/redigo/redis"
)

var redisPort string = "DB_PORT_6379"
var redisProto string = "tcp"

func Connect() (redis.Conn, error) {

	config := GetConfig()

	redisHost := config.RedisHost

	if redisHost == "" {
		os.Getenv(redisPort + "_TCP_ADDR")
	}

	// Redis test
	return redis.Dial(redisProto, redisHost)
}
