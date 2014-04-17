import logging

logging.basicConfig(level=logging.INFO)

from wsgiref import simple_server

from trajectory.service import app


# local
httpd = simple_server.make_server('127.0.0.1', 8888, app)
httpd.serve_forever()
