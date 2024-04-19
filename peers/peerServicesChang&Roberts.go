package main

import (
	"SDCCproject/algorithm"
	"SDCCproject/utils"
	"fmt"
)

func (PeerServiceHandler) NewLeaderCR(mex utils.Message, _ *utils.NodeINFO) error {
	if mex.MexID != currentNode.Id {
		currentNode.Leader = mex.MexID
	} else {
		return nil
	}

	currentNode.Participant = false

	/* Optional */
	for _, node := range currentNode.List.GetAllNodes() {
		currentNode.List.UpdateNode(node, mex.MexID)
	}

	go algorithm.WinnerMessage(currentNode, mex.MexID)

	return nil
}

func (PeerServiceHandler) ElectionMessageCR(mex utils.Message, _ *int) error {
	currID := currentNode.Id

	/* First message case */
	if mex.StartingMex == true {
		mex.StartingMex = false
		if mex.MexID > currID {
			currentNode.Participant = true

			go algorithm.ElectionChangAndRoberts(currentNode.List.GetNode(currID), mex.MexID)
		} else if mex.MexID < currID {
			currentNode.Participant = true
			mex.MexID = currID

			go algorithm.ElectionChangAndRoberts(currentNode.List.GetNode(currID), mex.MexID)
		}
	} else {
		if mex.MexID > currID {
			currentNode.Participant = true
			go algorithm.ElectionChangAndRoberts(currentNode.List.GetNode(currID), mex.MexID)
		} else if mex.MexID < currID && currentNode.Participant == false {
			currentNode.Participant = true
			mex.MexID = currID
			go algorithm.ElectionChangAndRoberts(currentNode.List.GetNode(currID), mex.MexID)
		} else if mex.MexID < currID && currentNode.Participant == true {
			fmt.Println("Algorithm error, message discarded")
			return nil // discard message
		} else if mex.MexID == currID {
			currentNode.Leader = currID
			currentNode.Participant = false

			go algorithm.WinnerMessage(currentNode, mex.MexID) // forward message to currentNode.next with leader: mex.MexID
		}
	}

	return nil
}
