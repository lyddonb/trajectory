import logging
import socketserver

from trajectory.stat import parse_stats


class TCPHandler(socketserver.BaseRequestHandler):

    def handle(self):
        self.data = self.request.recv(512).strip()

        self.data = self.data.decode('UTF-8')
        logging.info("Data: %s", self.data)

        if self.data[:5] == "STAT ":
            logging.info("Is stat")
            parse_stats(self.data[5:])


def start():
    HOST, PORT = "localhost", 1200

    # Create the server, binding to localhost on port 9999
    server = socketserver.TCPServer((HOST, PORT), TCPHandler)

    # Activate the server; this will keep running until you
    # interrupt the program with Ctrl-C
    server.serve_forever()


if __name__ == "__main__":
    start()
