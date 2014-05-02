import json
import logging

import falcon

from trajectory.machine import get_machines

from trajectory.stat import get_stats_for_machine

from trajectory.task import build_request_graph
from trajectory.task import get_task_for_key
from trajectory.task import iget_requests
from trajectory.task import iget_task_keys_for_request
from trajectory.task import to_tree


def write(req, resp, payload, as_json=True):
    host = req.get_header('origin')
    resp.set_header('Access-Control-Allow-Origin', host or "*")
    resp.set_header('Access-Control-Allow-Headers', "content-type, accept")
    resp.set_header('Access-Control-Allow-Methods',
                    'GET, POST, PUT, DELETE, OPTIONS')

    resp.status = falcon.HTTP_200

    if as_json:
        resp.body = json.dumps(payload)
    else:
        resp.body = payload


class MachineListResource:

    def on_get(self, req, resp):
        machines = get_machines(req.get_param('path'))

        write(req, resp, machines)


class StatListResource:

    def on_get(self, req, resp):
        machine_path = get_machines(req.get_param('path'))

        stats = get_stats_for_machine(machine_path)

        write(req, resp, stats)


class RequestListResource:

    def on_get(self, req, resp):
        write(req, resp, list(iget_requests()))


class RequestTasksListResource:

    def on_get(self, req, resp, request_id):
        write(req, resp, list(iget_task_keys_for_request(request_id)))


class TaskResource:

    def on_get(self, req, resp, key):
        write(req, resp, get_task_for_key(key))


class RequestTreeResource:

    def on_get(self, req, resp, request_id):
        write(req, resp, to_tree(build_request_graph(request_id)),
              as_json=False)


app = falcon.API()

app.add_route('/api/machines', MachineListResource())
app.add_route('/api/stats', StatListResource())
app.add_route('/api/requests', RequestListResource())
app.add_route('/api/requests/{request_id}', RequestTasksListResource())
app.add_route('/api/requests/{request_id}/tree', RequestTreeResource())
app.add_route('/api/task/{key}', TaskResource())
