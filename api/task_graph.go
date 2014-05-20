package api

import (
	"fmt"
	"strings"

	"github.com/lyddonb/trajectory/db"
)

type Node struct {
	TaskId    string
	Name      string
	ContextId string
	Children  map[string]*Node
	Keys      []string // Handles multiple runs of the same task
	IsParent  bool
}

func (a *TaskAPI) GetRequestTaskGraph(requestId string) (*Node, error) {
	taskKeys, e := a.ListRequestTaskKeys(requestId)

	if e != nil {
		return nil, e
	}

	// TODO: Hook in request info to get the url for this request.
	parent := BuildParentNode(requestId)

	if len(taskKeys) != 0 {
		a.ProcessTaskKeys(requestId, taskKeys, parent)
	}

	return parent, nil
}

func (a *TaskAPI) ProcessTaskKeys(requestId string, taskKeys map[string]int, parent *Node) {
	loadTaskChan := make(chan db.Task)

	requestNodes := make(map[string]*Node)
	requestNodes[requestId] = parent

	taskNodes := make(map[string]*Node)

	childrenMap := make(map[string]map[string]*Node)

	for taskKey, _ := range taskKeys {
		parentTaskId, node := ProcessChildNode(taskKey, requestNodes, childrenMap)

		if node.Name == "" {
			taskNodes[taskKey] = node
			go a.LoadTask(taskKey, loadTaskChan)
		}

		if parentTaskId == "" || strings.ToLower(parentTaskId) == "none" {
			parent.Children[taskKey] = node
		} else {
			parentNode, parentNodeExists := requestNodes[parentTaskId]

			if parentNodeExists {
				parentNode.Children[taskKey] = node
			} else {
				_, parentNodeExistsInChildren := childrenMap[parentTaskId]

				if parentNodeExistsInChildren {
					childrenMap[parentTaskId][taskKey] = node
				}
			}
		}
	}

	if len(taskNodes) > 0 {
		HandleTasks(loadTaskChan, taskNodes)
	}
}

func (a *TaskAPI) LoadTask(taskKey string, taskChannel chan<- db.Task) {
	// Load the task and pass it into the channel.
	task, _ := a.dal.GetTaskForKey(taskKey)

	// TODO: Handle the error
	taskChannel <- task
}

func HandleTasks(taskChannel <-chan db.Task, taskNodes map[string]*Node) {
	for {
		select {
		case task := <-taskChannel:
			node, isNode := taskNodes[task.Key()]

			if !isNode {
				fmt.Println("Error: couldn't find node?")
			}

			node.Name = task[db.URL]
			delete(taskNodes, task.Key())

			if len(taskNodes) == 0 {
				return
			}
		}
	}
}

func ProcessChildNode(taskKey string, nodes map[string]*Node,
	childrenMap map[string]map[string]*Node) (string, *Node) {

	parentTaskId, taskId, contextId := SplitTaskKey(taskKey)

	node, nodeExists := nodes[taskId]

	if nodeExists {
		node.Keys = append(node.Keys, taskKey)
	} else {
		children, childrenExist := childrenMap[taskId]

		if !childrenExist {
			children = make(map[string]*Node)
		}

		node = BuildChildNode(taskId, taskKey, contextId, children)
		nodes[taskId] = node

		// TODO: Load task info to get name
	}

	return parentTaskId, node
}

func BuildParentNode(requestId string) *Node {
	return &Node{
		requestId,
		"Parent Request",
		"",
		make(map[string]*Node),
		[]string{},
		true,
	}
}

func BuildChildNode(taskId, taskKey, contextId string,
	children map[string]*Node) *Node {
	return &Node{
		taskId,
		"",
		contextId,
		children,
		[]string{taskKey},
		false,
	}
}

func SplitTaskKey(taskKey string) (string, string, string) {
	splitKey := strings.Split(taskKey, ":")

	parentTaskId := splitKey[0]
	taskId := splitKey[1]
	var contextId string

	taskIdSplit := strings.Split(taskId, "|")

	if len(taskIdSplit) == 2 {
		taskId = taskIdSplit[0]
		contextId = taskIdSplit[1]
	}

	return parentTaskId, taskId, contextId
}
