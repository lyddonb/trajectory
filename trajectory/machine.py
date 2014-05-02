from collections import defaultdict

from trajectory.db import redisdb

from trajectory.stat import MACHINE_KEY


def get_machines(path=None):
    machines = redisdb.zrevrange(MACHINE_KEY, 0, 150)

    return_machines = {}

    # QUESTION: Maybe strip the . off the end of the path?

    for machine in machines:
        machine = machine.decode('UTF-8')

        parent = path

        if path:
            if not machine.startswith(path) or path == machine:
                continue

            match_path = path if path.endswith('.') else path + '.'

            machine = machine.replace(match_path, '', 1)
            print("***********************")
            print(machine)

            if not machine:
                continue

        name = machine.split('.')[0]

        if name in return_machines:
            continue

        return_machines[name] = make_machine(name, parent)

    return list(return_machines.values())


def make_machine(machine, parent):
    return {
        'machine': machine,
        'parent': parent
    }
