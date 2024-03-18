package utils

type Node struct {
	Id          int
	Address     string
	Leader      bool
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
	Leader      bool
	Participant bool
}

type LeaderStatus struct {
	NewLeaderID int
	OldLeaderID int
}

type Message struct {
	SkipCount int
	MexID     int
	CurrNode  NodeINFO
}
