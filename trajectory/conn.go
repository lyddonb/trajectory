package trajectory

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/garyburd/redigo/redis"
)

var port string = ":1200"
var Machine string = "Machine"

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

		statByte := buf[:n]

		// Check if stat type.
		if bytes.HasPrefix(statByte, []byte{'S', 'T', 'A', 'T', ' '}) {
			if redisClient == nil {
				fmt.Println("No redis client")
				fmt.Println(string(statByte))
			} else {
				handleStat(statByte, redisClient)
			}
		}
		// TODO: Handle Log
		// TODO: Handle Request Info

		if err != nil {
			// TODO: log the error.
			return
		}
	}
}

func handleStat(statByte []byte, redisClient redis.Conn) {
	// Kill the statBytePrefix
	timestamp := time.Now().UTC().Format("20060102150405")

	// TODO: Possibly, pass this back to a goroutine to apply to redis.
	for _, stat := range ParseStat(statByte[5:]) {
		redisClient.Send("HMSET", redis.Args{}.Add(stat.Key()).AddFlat(&stat)...)
		redisClient.Send("ZADD", Machine, timestamp, stat.Machine)
		redisClient.Send("ZADD", StatKeys, timestamp, stat.KeyWithOutStatType())
	}
	redisClient.Flush()

	v, _ := redisClient.Receive()
	fmt.Printf("Test %s\n", v)

	r, _ := redisClient.Do("KEYS", "Stat:*")
	fmt.Printf("Stats from redis %s.\n", r)

	s, _ := redisClient.Do("KEYS", "StatKeys:*")
	fmt.Printf("StatKeys from redis %s.\n", s)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
