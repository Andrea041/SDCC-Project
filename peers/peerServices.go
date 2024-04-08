package main

import (
	"SDCCproject/algorithm"
	"SDCCproject/utils"
	"fmt"
)

// PeerServiceHandler is the interface which exposes the RPC method ManageNode
type PeerServiceHandler struct{}

// UpdateList is a service to update node's list
func (PeerServiceHandler) UpdateList(node utils.Node, _ *utils.NodeINFO) error {
	currentNode.List.GetAllNodes()
	currentNode.List.AddNode(node)

	fmt.Println("New peer in system")

	return nil
}

// CheckLeaderStatus is a service that realize the ping to this node
func (PeerServiceHandler) CheckLeaderStatus(callerNode utils.Node, _ *utils.NodeINFO) error {
	fmt.Printf("Ping received from process %d\n", callerNode.Id)

	return nil
}

// ElectionMessageBULLY is a service that reply "OK" in rep to the caller node
func (PeerServiceHandler) ElectionMessageBULLY(nodeCaller utils.NodeINFO, rep *string) error {
	fmt.Printf("Election message from Peer with id: %d\n", nodeCaller.Id)
	*rep = "OK"

	go algorithm.ElectionBully(currentNode)

	return nil
}

func (PeerServiceHandler) NewLeaderBULLY(leaderNode utils.NodeINFO, _ *utils.NodeINFO) error {
	currentNode.Leader = leaderNode.Id

	/* Optional */
	for _, node := range leaderNode.List.GetAllNodes() {
		currentNode.List.UpdateNode(node, leaderNode.Id)
	}

	return nil
}

func (PeerServiceHandler) NewLeaderCR(mex utils.Message, _ *utils.NodeINFO) error {
	if mex.MexID != currentNode.Id {
		currentNode.Leader = mex.MexID
	} else {
		return nil
	}

	/* Optional */
	for _, node := range currentNode.List.GetAllNodes() {
		currentNode.List.UpdateNode(node, mex.MexID)
	}

	go algorithm.WinnerMessage(currentNode, mex.MexID)

	return nil
}

func (PeerServiceHandler) ElectionMessageCR(mex utils.Message, _ *int) error {
	currID := currentNode.Id

	if mex.MexID > currID {
		go algorithm.ElectionChangAndRoberts(currentNode.List.GetNode(currID), mex.MexID)
	} else if mex.MexID < currID {
		mex.MexID = currID
		go algorithm.ElectionChangAndRoberts(currentNode.List.GetNode(currID), mex.MexID)
	} else if mex.MexID == currID {
		fmt.Println("Forwarding winner message")
		currentNode.Leader = currID

		go algorithm.WinnerMessage(currentNode, mex.MexID)
	}

	return nil
}
