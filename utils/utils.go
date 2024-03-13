package utils

type Node struct {
	Id      int
	Address string
	Leader  bool
}

type PeerAddr struct {
	PeerAddress string
}

type NodeList struct {
	Nodes []Node
}

type NodeINFO struct {
	Id     int
	List   NodeList
	Leader bool
}

func (nl *NodeList) AddNode(node Node) {
	nl.Nodes = append(nl.Nodes, node)
}

func (nl *NodeList) GetAllNodes() []Node {
	return nl.Nodes
}
