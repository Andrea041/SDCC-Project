package main

import (
	"SDCCproject/algorithm"
	"SDCCproject/utils"
)

// ElectionMessageBULLY is a service that reply "OK" in rep to the caller node
func (PeerServiceHandler) ElectionMessageBULLY(_ utils.NodeINFO, rep *string) error {
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
