import logging
import socketserver

from trajectory.stat import parse_stats
from trajectory.task import parse_tasks

STAT = "STAT"
TASK = "TASK"
PADDING = " "
KEY_SIZE = 5


HANDLERS = {
    STAT: parse_stats,
    TASK: parse_tasks
}


class TCPHandler(socketserver.BaseRequestHandler):

    def handle(self):
        self.data = self.request.recv(1024).strip()

        self.data = self.data.decode('UTF-8')
        logging.info("Data: %s", self.data)

        prefix = self.data[:KEY_SIZE]
        key = prefix.rstrip(PADDING)

        if key not in HANDLERS:
            logging.info("Unhandled handler type.")
            return

        HANDLERS[key](self.data[KEY_SIZE:])


def start():
    HOST, PORT = "localhost", 1200

    # Create the server, binding to localhost on port 9999
    server = socketserver.TCPServer((HOST, PORT), TCPHandler)

    # Activate the server; this will keep running until you
    # interrupt the program with Ctrl-C
    server.serve_forever()


if __name__ == "__main__":
    start()
