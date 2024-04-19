package utils

type Node struct {
	Id          int
	Address     string
	Leader      int
	Participant bool
}

type PeerAddr struct {
	PeerAddress string
}

type NodeList struct {
	Nodes []Node
}

type NodeINFO struct {
	Id          int
	Address     string
	List        NodeList
	Leader      int
	Participant bool
}

type Message struct {
	SkipCount   int
	MexID       int
	CurrNode    NodeINFO
	StartingMex bool
}

type Address struct {
	Address string `json:"address"`
	Port    string `json:"port"`
}

type Alg struct {
	Chang string `json:"ChangAndRobert"`
	Bully string `json:"Bully"`
}

type Configuration struct {
	ServiceRegistry Address `json:"service_registry"`
	Peer            Address `json:"peer"`
	Algorithm       Alg     `json:"algorithm"`
}
