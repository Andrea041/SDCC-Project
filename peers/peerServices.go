package main

import (
	"SDCCproject/utils"
	"fmt"
)

// PeerServiceHandler is the interface which exposes the RPC method ManageNode
type PeerServiceHandler struct{}

// UpdateList is a service to update node's list
func (PeerServiceHandler) UpdateList(node utils.Node, _ *utils.NodeINFO) error {
	currentNode.List.GetAllNodes()
	currentNode.List.AddNode(node)

	return nil
}

// CheckLeaderStatus is a service that realize the ping to this node
func (PeerServiceHandler) CheckLeaderStatus(callerNode utils.Node, _ *utils.NodeINFO) error {
	fmt.Printf("Ping received from process %d\n", callerNode.Id)

	return nil
}
