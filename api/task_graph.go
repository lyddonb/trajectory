package api

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/GaryBoone/GoStats/stats"
	"github.com/lyddonb/trajectory/db"
)

type Node struct {
	TaskId      string  `json:"task_id"`
	Key         string  `key`
	Name        string  `json:"name"`
	ContextId   string  `json:"context_id"`
	Children    []*Node `json:"children"`
	childrenMap map[string]*Node
	Keys        []string `json:"keys"` // Handles multiple runs of the same task
	IsParent    bool     `json:"is_parent"`
	Status      int      `json:"status"`
	Latency     float64  `json:"latency"`
	RunTime     float64  `json:"run_time"`
	task        *db.Task
	Stats       map[string]*Stats
}

type Stats struct {
	Count   int     `json:"count"`
	Min     float64 `json:"min"`
	Max     float64 `json:"max"`
	Average float64 `json:"average"`
	Median  float64 `json:"median"`
	StdDev  float64 `json:"std_dev"`
}

//func (n *Node) MarshalJSON() ([]byte, error) {
//}

func convertAndUpdateStat(value string, nodeStat *stats.Stats) {
	runValue, err := strconv.ParseFloat(value, 64)

	if err != nil {
		fmt.Println(err)
		return
	}

	nodeStat.Update(runValue)
}

func BuildStats(parent *Node) {
	nodeStats := make(map[string]*stats.Stats)
	props := []string{"end", "execution_count", "gae_latency_seconds",
		"ran", "retry_count", "run_time", "status_code", "task_eta"}

	for _, prop := range props {
		nodeStats[prop] = &stats.Stats{}
	}

	for _, child := range parent.Children {
		task := *child.task
		BuildStats(child)
		for _, prop := range props {
			convertAndUpdateStat(task[prop], nodeStats[prop])
		}
	}

	for key, value := range nodeStats {
		runStat := &Stats{0, 0, 0, 0, 0, 0}
		runStat.Count = value.Count()
		runStat.Min = value.Min()
		runStat.Max = value.Max()
		runStat.Average = value.Mean()
		stddev := value.SampleStandardDeviation()

		if !math.IsNaN(stddev) {
			runStat.StdDev = stddev
		}

		parent.Stats[key] = runStat
	}

}

func (a *TaskAPI) GetRequestTaskGraph(requestId string) (*Node, error) {
	taskKeys, e := a.ListRequestTaskKeysAsSet(requestId)

	if e != nil {
		return nil, e
	}

	// TODO: Hook in request info to get the url for this request.
	parent := BuildParentNode(requestId)

	if len(taskKeys) != 0 {
		parent = a.ProcessTaskKeys(requestId, taskKeys, parent)
	}

	BuildStats(parent)

	return parent, nil
}

func (a *TaskAPI) ProcessTaskKeys(requestId string, taskKeys map[string]int,
	parent *Node) *Node {
	//loadTaskChan := make(chan db.Task)

	requestNodes := make(map[string]*Node)
	requestNodes[requestId] = parent

	taskNodes := make(map[string]*Node)

	childrenMap := make(map[string]map[string]*Node)

	for taskKey, _ := range taskKeys {
		parentTaskId, node := ProcessChildNode(taskKey, requestNodes, childrenMap)

		if node.Name == "" {
			taskNodes[taskKey] = node
			a.SetTaskInfo(taskKey, node)
			//go a.LoadTask(taskKey, loadTaskChan)
		}

		if parentTaskId == "" || strings.ToLower(parentTaskId) == "none" {
			parent.childrenMap[taskKey] = node
		} else {
			parentNode, parentNodeExists := requestNodes[parentTaskId]
			childItems, parentNodeExistsInChildren := childrenMap[parentTaskId]

			if parentNodeExistsInChildren || !parentNodeExists {
				if !parentNodeExistsInChildren {
					childrenMap[parentTaskId] = make(map[string]*Node)
				}
				childrenMap[parentTaskId][taskKey] = node
			}

			if parentNodeExists {
				if parentNodeExistsInChildren {
					parentNode.childrenMap = childItems
					delete(childrenMap, parentTaskId)
				} else {
					parentNode.childrenMap[taskKey] = node
				}
			}
		}
	}

	for _, node := range requestNodes {
		node.Children = make([]*Node, len(node.childrenMap))
		index := 0
		for _, child := range node.childrenMap {
			node.Children[index] = child
			index++
		}
	}

	//if len(taskNodes) > 0 {
	//HandleTasks(loadTaskChan, taskNodes)
	//}
	return parent
}

func GetStatFromTask(task db.Task, prop string) float64 {
	stat, err := strconv.ParseFloat(task[prop], 64)

	if err != nil {
		fmt.Println("Error parsing ", prop, err)
		return 0
	}

	return stat
}

func (a *TaskAPI) SetTaskInfo(taskKey string, node *Node) {
	// Load the task and pass it into the channel.
	task, err := a.dal.GetTaskForKey(taskKey)

	if err != nil {
		fmt.Println(err)
	}

	node.Name = task[db.URL]
	node.task = &task
	node.Latency = GetStatFromTask(task, db.LATENCY)
	node.RunTime = GetStatFromTask(task, db.RUN_TIME)

	status := 0
	if val, ok := task[db.STATUS_CODE]; ok {
		if val == "200" {
			status = 1
		} else {
			status = 2
		}
	}

	node.Status = status
}

func (a *TaskAPI) LoadTask(taskKey string, taskChannel chan<- db.Task) {
	// Load the task and pass it into the channel.
	task, err := a.dal.GetTaskForKey(taskKey)

	if task.Key() != taskKey {
		fmt.Println(task.Key())
		fmt.Println(taskKey)
		fmt.Println(task)
		//task[db.TASK_ID]
	}

	if err != nil {
		fmt.Println(err)
	}

	// TODO: Handle the error
	taskChannel <- task
}

func HandleTasks(taskChannel <-chan db.Task, taskNodes map[string]*Node) {
	for {
		select {
		case task := <-taskChannel:
			node, isNode := taskNodes[task.Key()]

			if !isNode {
				fmt.Println("Error: couldn't find node? %s", task.Key())
				fmt.Println(task)
				//return
			} else {
				node.Children = make([]*Node, len(node.childrenMap))
				index := 0
				for _, child := range node.childrenMap {
					node.Children[index] = child
					index++
				}

				node.Name = task[db.URL]
			}

			delete(taskNodes, task.Key())

			if len(taskNodes) == 0 {
				fmt.Println("finished processing")
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
		"",
		"Parent Request",
		"",
		make([]*Node, 0),
		make(map[string]*Node),
		[]string{},
		true,
		0,
		0,
		0,
		nil,
		make(map[string]*Stats),
	}
}

func BuildChildNode(taskId, taskKey, contextId string,
	children map[string]*Node) *Node {
	return &Node{
		taskId,
		taskKey,
		"",
		contextId,
		make([]*Node, 0),
		children,
		[]string{taskKey},
		false,
		0,
		0,
		0,
		nil,
		make(map[string]*Stats),
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
