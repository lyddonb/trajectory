import json
import logging
import time

from collections import defaultdict

from trajectory.db import redisdb


TASK_PAYLOAD_ID = "id"
REQUEST_INFO = "request_info"

TASK_ID = "task_id"
PARENT_TASK_ID = "parent_task_id"
PARENT_REQUEST_ID = "parent_request_id"
REQUEST_ID = "request_id"
SERVER_INFO = "server_info"
URL = 'url'

PARENT_REQUESTS = "ParentRequests"


class Task:

    def __init__(self, **kwargs):

        for k, v in kwargs.items():
            setattr(self, k, v)

    def to_dict(self):
        return self.__dict__

    @classmethod
    def from_dict(cls, dct):
        return cls(**dct)

    @property
    def key(self):
        return "%s:%s:%s" % (self.parent_task_id, self.task_id,
                             self.request_id)


def parse_tasks(task_input):
    task = build_task(task_input)

    save_task(task)


def build_task(task_input):
    task_payload = json.loads(task_input)

    _id = task_payload.pop(TASK_PAYLOAD_ID)

    split_id = _id.split(":")

    if len(split_id) == 3:
        parent_task_id = split_id[1]
        parent_request_id = split_id[0]
        task_id = split_id[2]
    else:
        parent_task_id = None
        parent_request_id = split_id[0]
        task_id = split_id[1]

    server_info, request_id = task_payload.pop(REQUEST_INFO).split("#")

    logging.info(task_payload)

    return Task(task_id=task_id,
                parent_task_id=parent_task_id,
                parent_request_id=parent_request_id,
                request_id=request_id,
                server_info=server_info,
                **task_payload)


def save_task(task):
    timestamp = time.time()

    pipe = redisdb.pipeline()
    pipe.hmset(task.key, task.to_dict())
    pipe.zadd(task.task_id, timestamp, task.key)
    pipe.zadd(task.parent_request_id, timestamp, task.key)
    pipe.zadd(PARENT_REQUESTS, timestamp, task.parent_request_id)
    pipe.execute()


def iget_requests():
    return (request_id.decode('UTF-8')
            for request_id in redisdb.zrevrange(PARENT_REQUESTS, 0, 150))


def iget_task_keys_for_request(parent_request_id):
    return (task_key.decode('UTF-8')
            for task_key in redisdb.zrevrange(parent_request_id, 0, 150))


class Node:

    def __init__(self, task_id=None, keys=None, name=None, children=None,
                 is_parent=False):
        self.task_id = task_id
        self.keys = keys or set()
        self.name = name
        self.children = children or []
        self.is_parent = is_parent

    def to_dict(self):
        return {
            'name': self.name,
            'task_id': self.task_id,
            'children': self.children,
            'keys': list(self.keys),
            'is_parent': self.is_parent
        }


def build_request_graph(parent_request_id):
    # TODO: Get request info for url to stick as name
    parent = Node(task_id=parent_request_id, is_parent=True)
    nodes = {
        parent_request_id: parent
    }
    children = defaultdict(list)

    # Might want to make the name the task_key.

    # TODO: Use asyncid to load the task info and stick it in the node info.
    # TODO: Group task ids together

    for task_key in iget_task_keys_for_request(parent_request_id):
        parent_task_id, task_id, _ = task_key.split(':')

        task_id = task_id.split('|')[0]

        if task_id in nodes:
            node = nodes[task_id]
            node.keys.add(task_key)
        else:
            node = Node(task_id=task_id, keys=set([task_key]))

            # Check for children and add them.
            if task_id in children:
                node.children = children.pop(task_id)

            nodes[task_id] = node

            task_info = get_task_for_key(task_key)
            if task_info:
                node.name = task_info[URL]

        if not parent_task_id or parent_task_id.lower() == "none":
            parent.children.append(node)
        elif parent_task_id not in nodes:
            if node not in children[parent_task_id]:
                children[parent_task_id].append(node)
        elif node not in nodes[parent_task_id].children:
            nodes[parent_task_id].children.append(node)

    return parent


def to_tree(parent):
    return json.dumps(parent, cls=ObjectEncoder)


def get_task_for_key(key):
    task_dict = {}

    for k, v in redisdb.hgetall(key).items():
        task_dict[k.decode('UTF-8')] = v.decode('UTF-8')

    return task_dict


class ObjectEncoder(json.JSONEncoder):

    def default(self, obj):

        if hasattr(obj, 'to_dict'):
            return obj.to_dict()
