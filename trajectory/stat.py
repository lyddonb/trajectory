import json
import logging
import time

from trajectory.db import redisdb


# TODO: Add timestamp to stat


MACHINE_KEY = "Machine"
STAT_KEY = "StatKeys"

STAT_KEY_PREFIX = "Stat:"

STAT_STATS = 'stats'
STAT_ID = 'id'
STAT_MACHINE = 'machine'
STAT_VALUE = 'value'
STAT_TYPE = 'type'


class Stat():

    def __init__(self, id=None, machine=None, value=None, type=None):
        self.id = id
        self.machine = machine
        self.value = value
        self.type = type

    @property
    def key(self):
        return STAT_KEY_PREFIX + self.id

    @property
    def full_key(self):
        return "%s$%s" % (self.machine, ".".join(self.id.split(".")[:-1]))

    def to_dict(self):
        return {
            STAT_ID: self.id,
            STAT_MACHINE: self.machine,
            STAT_TYPE: self.type,
            STAT_TYPE: self.value
        }

    @classmethod
    def from_dict(cls, dct):
        return cls(**dct)


def parse_stats(stat_input):
    stats = ibuild_stats(stat_input)

    # TODO: Move this to asyncio
    save_stats(stats)


def save_stats(stats):
    timestamp = time.time()

    for stat in stats:
        pipe = redisdb.pipeline()
        pipe.hmset(stat.key, stat.to_dict())
        pipe.zadd(MACHINE_KEY, timestamp, stat.machine)
        pipe.zadd(STAT_KEY, timestamp, stat.full_key)
        pipe.execute()


# TODO: Convert to generator.
def ibuild_stats(stat_input):
    logging.info("Parsing %s to json", stat_input)
    stats_paylaod = json.loads(stat_input)

    logging.info(stats_paylaod)

    stats = (Stat("%s.%s" % (stats_paylaod[STAT_ID], stat_key),
                  stats_paylaod[STAT_MACHINE], *value.split('|'))
             for stat_key, value in stats_paylaod[STAT_STATS].items())

    return stats


def get_stats_for_machine(path=None):
    # If no path return all stats.

    # TODO: Add path filter.
    return [make_stat(stat) for stat in redisdb.zrevrange(STAT_KEY, 0, 150)]


def make_stat(stat):
    stat = stat.decode('UTF-8')

    # Kinda GAE specific. Tweak later.
    split_stat = stat.split('$')
    request_id, url = split_stat[1].split('.')

    return {
        'stat_key': split_stat[1],
        'request_id': request_id,
        'url': url,
        'parent': split_stat[0]
    }
