package utils

func (nl *NodeList) AddNode(node Node) {
	nl.Nodes = append(nl.Nodes, node)
}

func (nl *NodeList) GetAllNodes() []Node {
	return nl.Nodes
}

func (nl *NodeList) UpdateNode(node Node, update bool) {
	for i, nodeToUpdate := range nl.Nodes {
		if nodeToUpdate.Id == node.Id {
			(nl.Nodes)[i].Leader = update
			return
		}
	}
}
