package utils

type Node struct {
	Id      int
	Address string
	Leader  int
}

type PeerAddr struct {
	PeerAddress string
}

type NodeList struct {
	Nodes []Node
}

type NodeINFO struct {
	Id      int
	Address string
	List    NodeList
	Leader  int
}

type Message struct {
	SkipCount int
	MexID     int
	CurrNode  NodeINFO
}

type Address struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type Configuration struct {
	ServiceRegistry Address `json:"service_registry"`
	Peer            Address `json:"peer"`
}
