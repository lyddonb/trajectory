package api

import "fmt"

type NodeLink struct {
	Node   *Node
	Parent *NodeLink
}

func (a *TaskAPI) GetFullGraph(requestId string, endId string) ([]*Node, error) {
	graph, e := a.GetRequestTaskGraph(requestId)

	if e != nil {
		return nil, e
	}

	nodeLink := &NodeLink{graph, nil}

	fmt.Println("endId = ", endId)
	node := FindNode(endId, nodeLink)
	fmt.Println("FindNode Results = ", node)

	var slice []*Node
	slice = append(slice, node.Node)

	chain := ParentList(slice, node)
	fmt.Println("Chain: ", chain)

	// TODO: Return the list of nodes

	return chain, nil
}

func ParentList(slice []*Node, node *NodeLink) []*Node {
	//parentLink := NodeLink{node.Parent.Node, node.Parent.Parent}

	if node.Parent == nil {
		return slice
	} else {
		slice = append(slice, node.Parent.Node)
		//fmt.Println(slice)
		result := ParentList(slice, node.Parent)
		if result != nil {
			return result
		}
	}
	return slice
}

func FindNode(nodeId string, parent *NodeLink) *NodeLink {
	for _, child := range parent.Node.Children {
		childLink := NodeLink{child, parent}

		task := *child.task
		fmt.Println("$$$", task["task_id"])

		if task["task_id"] == nodeId {
			return &childLink
		} else {
			fmt.Println("***", task["task_id"])
			result := FindNode(nodeId, &childLink)
			if result != nil {
				return result
			}
		}
	}

	return nil
}
