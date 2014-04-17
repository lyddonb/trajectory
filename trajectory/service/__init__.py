import json

import falcon

from trajectory.db import redisdb

from trajectory.stat import MACHINE_KEY
from trajectory.stat import STAT_KEY


class MachineListResource:

    def on_get(self, req, resp):
        machines = redisdb.zrevrange(MACHINE_KEY, 0, 150)

        resp.status = falcon.HTTP_200
        resp.body = json.dumps(
            [machine.decode('UTF-8') for machine in machines])


class StatListResource:

    def on_get(self, req, resp):
        stats = redisdb.zrevrange(STAT_KEY, 0, 150)

        resp.status = falcon.HTTP_200
        resp.body = json.dumps(
            [stat.decode('UTF-8') for stat in stats])


app = falcon.API()

app.add_route('/api/machines', MachineListResource())
app.add_route('/api/stats', StatListResource())
