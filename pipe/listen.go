package pipe

import (
	"fmt"
	"net"
	"os"
)

const connectionSize = 1024

type Pipeline interface {
	Handler(net.Conn)
	Error(error)
	Open() bool
	Parse([]byte, string)
}

func MakeConnection(connection, address string) net.Listener {
	fmt.Printf("Opening a %s address on %s.\n", connection, address)

	listener, err := net.Listen(connection, address)
	checkErrorAndExit(err)

	return listener
}

// Open a tcp or udp connection on the passed in address that listens for input.
func Listen(listener net.Listener, pipeline Pipeline) {

	defer listener.Close()

	for pipeline.Open() {
		socketConn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error %s trying to listen on %s.", err, listener)
			go pipeline.Error(err)
			continue
		}

		go pipeline.Handler(socketConn)
	}

	fmt.Printf("Closing listener.\n")
}

func handleClient(socketConn net.Conn, pipeline Pipeline) {
	var buf = make([]byte, connectionSize)

	defer socketConn.Close()

	for pipeline.Open() {
		count, err := socketConn.Read(buf[0:])

		if err != nil {
			fmt.Println("Error %s trying to accept data.", err)
			pipeline.Error(err)
			return
		}

		// TODO: Have parse data pass to a channel to apply to the database,
		// this can be quite large as redis can take in a ton of connections.
		go pipeline.Parse(buf[:count], socketConn.RemoteAddr().String())
	}
}

// Check for an error and exit if it's hit.
func checkErrorAndExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
