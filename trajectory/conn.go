package trajectory

import (
	"fmt"
	"net"
	"os"

	"github.com/garyburd/redigo/redis"
)

var port string = ":1200"

func Listen(redisClient redis.Conn) {
	connection := "tcp"
	fmt.Printf("Opening a %s port on %s.\n", connection, port)

	listener, err := net.Listen(connection, port)
	checkError(err)

	for {
		socketConn, err := listener.Accept()

		if err != err {
			// TODO: log the error.
			continue
		}

		go handleClient(socketConn, redisClient)
	}
}

func handleClient(socketConn net.Conn, redisClient redis.Conn) {
	var buf [512]byte

	for {
		n, err := socketConn.Read(buf[0:])

		if err != nil {
			// TODO: log the error.
			return
		}

		for _, stat := range ParseStat(buf[:n]) {
			redisClient.Send("HSET", stat.RedisKey(), stat.Id, stat.ToRedis())
		}
		redisClient.Flush()

		v, _ := redisClient.Receive()
		fmt.Printf("Test %s\n", v)

		r, _ := redisClient.Do("KEYS", "Stat:*")
		fmt.Printf("Response from redis %s.\n", r)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
