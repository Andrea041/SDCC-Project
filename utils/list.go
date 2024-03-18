package utils

func (nl *NodeList) AddNode(node Node) {
	nl.Nodes = append(nl.Nodes, node)
}

func (nl *NodeList) GetAllNodes() []Node {
	return nl.Nodes
}

func (nl *NodeList) GetNode(nodeId int) NodeINFO {
	for _, nodeWanted := range nl.Nodes {
		if nodeWanted.Id == nodeId {
			return NodeINFO{nodeWanted.Id, nodeWanted.Address, *nl, nodeWanted.Leader, nodeWanted.Participant}
		}
	}
	return NodeINFO{}
}

func (nl *NodeList) UpdateNode(node Node, update bool) {
	for i, nodeToUpdate := range nl.Nodes {
		if nodeToUpdate.Id == node.Id {
			(nl.Nodes)[i].Leader = update
			return
		}
	}
}
