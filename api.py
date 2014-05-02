import logging

logging.basicConfig(level=logging.INFO)

from wsgiref import simple_server

from trajectory.service import app


def start():
    # local
    server = '127.0.0.1'
    port = 8888
    print("serving at %s:%s" % (server, port))
    httpd = simple_server.make_server(server, port, app)
    httpd.serve_forever()


if __name__ == "__main__":
    start()
