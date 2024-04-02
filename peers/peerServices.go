package main

import (
	"SDCCproject/algorithm"
	"SDCCproject/utils"
	"time"

	"fmt"
	"log"
)

// PeerServiceHandler is the interface which exposes the RPC method ManageNode
type PeerServiceHandler struct{}

// UpdateList is a service to update node's list
func (PeerServiceHandler) UpdateList(node utils.Node, _ *utils.NodeINFO) error {
	currentNode.List.GetAllNodes()
	currentNode.List.AddNode(node)

	/*shuffledList := make([]utils.Node, len(currentNode.List.Nodes))
	perm := rand.Perm(len(currentNode.List.Nodes))
	for i, v := range perm {
		shuffledList[i] = currentNode.List.Nodes[v]
	}
	currentNode.List.Nodes = shuffledList
	*/

	fmt.Printf("New peer in system, list updated %s\n", currentNode.List.Nodes)

	return nil
}

// CheckLeaderStatus is a service that realize the ping to this node
func (PeerServiceHandler) CheckLeaderStatus(_ utils.Node, _ *utils.NodeINFO) error {
	fmt.Printf("Hi i'm the leader with address: %s, and id: %d. I'm still active!\n", currentNode.Address, currentNode.Id)

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
	nextNode := (mex.CurrNode.Id + mex.SkipCount) % len(currentNode.List.Nodes)

	if currentNode.Leader != mex.MexID {
		currentNode.Leader = mex.MexID
	} else {
		return nil
	}

	/* Optional */
	for _, node := range currentNode.List.GetAllNodes() {
		currentNode.List.UpdateNode(node, mex.MexID)
	}

	go algorithm.WinnerMessage(currentNode.List.GetNode(nextNode), mex.MexID)

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
		info := utils.Message{
			SkipCount: 1,
			MexID:     mex.MexID,
			CurrNode:  currentNode.List.GetNode(currID),
		}

		currentNode.Leader = mex.MexID

		startIndex := currentNode.List.GetIndex(currentNode.Id)
		nextNode := currentNode.List.Nodes[(startIndex+1)%len(currentNode.List.Nodes)]

		peer, err := utils.DialTimeout("tcp", currentNode.List.GetNode(nextNode.Id).Address, 5*time.Second)
		if err != nil {
			skip := startIndex
			i := 1
			for {
				i++
				pass := (startIndex + i) % len(currentNode.List.Nodes)
				if pass == skip-1 {
					return nil
				}

				peer, err = utils.DialTimeout("tcp", currentNode.List.GetNode(currentNode.List.Nodes[pass].Id).Address, 5*time.Second)
				info.SkipCount = i
				if err != nil {
					continue
				} else {
					break
				}
			}
		}

		err = peer.Call("PeerServiceHandler.NewLeaderCR", info, nil)
		if err != nil {
			log.Printf("New leader update error: %v", err)
		}

		err = peer.Close()
		if err != nil {
			log.Fatalf("Closing connection error: %v\n", err)
		}
	}

	return nil
}
