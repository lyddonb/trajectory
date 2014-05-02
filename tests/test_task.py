from mock import patch

from trajectory.task import build_request_graph


@patch('trajectory.task.iget_task_keys_for_request')
def test_single_node(get_tasks):
    """Ensure a single result for the parent returns a single node as the
    parent.
    """
    parent_request_id = 'requestid'
    key = 'None:taskid:foo'

    get_tasks.return_value = iter([key])

    parent = build_request_graph(parent_request_id)

    assert parent.task_id == 'requestid'
    assert len(parent.children) == 1
    assert parent.is_parent

    child = parent.children[0]

    assert child.task_id == "taskid"
    assert key in child.keys
    assert not child.is_parent
    assert len(child.children) == 0

    get_tasks.assert_called_once_with(parent_request_id)


@patch('trajectory.task.iget_task_keys_for_request')
def test_parent_child_in_order(get_tasks):
    """Ensure a parent child relationship in the hierarchical order returns a
    parent node with a single child.
    """
    parent_request_id = 'requestid'
    task_id = 'taskid'
    parent_key = 'None:%s:foo' % (task_id,)
    child_key = '%s:child1taskid:foo' % (task_id,)

    get_tasks.return_value = iter([parent_key, child_key])

    parent = build_request_graph(parent_request_id)

    assert parent.key == parent_key
    assert len(parent.children) == 1

    child = parent.children[0]

    assert child.key == child_key
    assert child.children == []


@patch('trajectory.task.iget_task_keys_for_request')
def test_parent_child_out_of_order(get_tasks):
    """Ensure a parent child relationship not in the hierarchical order returns
    a parent node with a single child.
    """
    parent_request_id = 'requestid'
    task_id = 'taskid'
    child_key = '%s:child1taskid:foo' % (task_id,)
    parent_key = 'None:%s:foo' % (task_id,)

    get_tasks.return_value = iter([child_key, parent_key])

    parent = build_request_graph(parent_request_id)

    assert parent.task_id == 'requestid'
    assert len(parent.children) == 1
    assert parent.is_parent

    child = parent.children[0]

    assert child.task_id == "taskid"
    assert parent_key in child.keys
    assert not child.is_parent
    assert len(child.children) == 1

    child2 = child.children[0]

    assert child2.task_id == "child1taskid"
    assert child_key in child2.keys
    assert not child.is_parent
    assert len(child2.children) == 0



@patch('trajectory.task.iget_task_keys_for_request')
def test_parent_child_child_out_of_order(get_tasks):
    """Ensure 3 deep parent child relationship not in the hierarchical order
    returns a parent node with a child and it also has a child.
    """
    parent_request_id = 'requestid'
    task_id = 'taskid'
    child1_task_id = 'child1taskid'
    child1_key = '%s:%s:foo' % (task_id, child1_task_id,)
    child2_key = '%s:child2taskid:foo' % (child1_task_id,)
    parent_key = 'None:%s:foo' % (task_id,)

    get_tasks.return_value = iter([child1_key, child2_key, parent_key])

    parent = build_request_graph(parent_request_id)

    assert parent.key == parent_key
    assert len(parent.children) == 1

    child1 = parent.children[0]

    assert child1.key == child1_key
    assert len(child1.children) == 1

    child2 = child1.children[0]

    assert child2.key == child2_key
    assert child2.children == []


@patch('trajectory.task.iget_task_keys_for_request')
def test_parent_2_children_out_of_order(get_tasks):
    """Ensure a parent with 2 children returns the parent with the 2 children.
    """
    parent_request_id = 'requestid'
    task_id = 'taskid'
    child1_task_id = 'child1taskid'
    child2_task_id = 'child2taskid'
    child1_key = '%s:%s:foo' % (task_id, child1_task_id,)
    child2_key = '%s:%s:foo' % (task_id, child2_task_id,)
    parent_key = 'None:%s:foo' % (task_id,)

    get_tasks.return_value = iter([child1_key, child2_key, parent_key])

    parent = build_request_graph(parent_request_id)

    assert parent.key == parent_key
    assert len(parent.children) == 2

    child1 = parent.children[0]

    assert child1.key == child1_key
    assert child1.children == []

    child2 = parent.children[1]

    assert child2.key == child2_key
    assert child2.children == []
