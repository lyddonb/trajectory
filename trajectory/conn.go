package trajectory

import (
	"fmt"
	"net"
	"os"

	"github.com/garyburd/redigo/redis"
)

/*
{
    "request": "requestid.path.appid.version.module.instanceid",
    "items": {
        "cpu_usage": "0.0|g",
        "memory_usage": "0.0|g",
        "exec_time": "0|ms",
        "overhead": "0|g",
        "system_GetSystemStats_offset": "6|g",
        "end_memory": "0.0|g",
        "end_cpu": "0.0|g"
    }
}
*/

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

		s := string(buf[:n])

		fmt.Printf("Sent item: %s\n", s)

		//redisClient.Do("SET", "hello", "world")
		//r, _ := redisClient.Do("GET", "hello")
		//fmt.Printf("Response from redis %s.\n", r)

		//_, err2 := conn.Write(buf[0:n])

		//if err2 != nil {
		//return
		//}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
