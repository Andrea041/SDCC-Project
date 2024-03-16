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

type LeaderStatus struct {
	NewLeaderID int
	OldLeaderID int
}
