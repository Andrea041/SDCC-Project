package algorithm

import (
	"SDCCproject/utils"

	"fmt"
	"log"
	"net/rpc"
	"time"
)

func WinnerMessage(currentNode utils.NodeINFO, leader int) {
	info := utils.Message{
		SkipCount: 1,
		MexID:     leader,
		CurrNode:  currentNode,
	}

	startIndex := currentNode.List.GetIndex(currentNode.Id)
	nextNode := currentNode.List.Nodes[(startIndex+1)%len(currentNode.List.Nodes)]

	peer, err := utils.DialTimeout("tcp", currentNode.List.GetNode(nextNode.Id).Address, 5*time.Second)
	if err != nil {
		skip := startIndex
		i := 0
		for {
			i++
			pass := (startIndex + i) % len(currentNode.List.Nodes)
			if pass == skip-1 {
				return
			}

			peer, err = utils.DialTimeout("tcp", currentNode.List.GetNode(currentNode.List.Nodes[pass].Id).Address,
				5*time.Second)
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
		log.Fatal("Leader update error", err)
	}

	err = peer.Close()
	if err != nil {
		log.Fatal("Closing connection error: ", err)
	}
}

func ElectionChangAndRoberts(currentNode utils.NodeINFO, mexReply int) {
	var info utils.Message
	info = utils.Message{SkipCount: 1, MexID: mexReply, CurrNode: currentNode}

	startIndex := currentNode.List.GetIndex(currentNode.Id)
	nextNode := currentNode.List.Nodes[(startIndex+1)%len(currentNode.List.Nodes)]

	peer, err := utils.DialTimeout("tcp", currentNode.List.GetNode(nextNode.Id).Address, 5*time.Second)
	if err != nil {
		/* Calculate skip for node inactivity */
		skip := startIndex
		i := 0
		for {
			i++
			pass := (startIndex + i) % len(currentNode.List.Nodes)
			if pass == skip-1 {
				info.MexID = currentNode.Id
				peer, err = utils.DialTimeout("tcp", currentNode.List.GetNode(currentNode.Id).Address, 5*time.Second)

				err = peer.Call("PeerServiceHandler.NewLeaderCR", info, nil)
				if err != nil {
					log.Printf("Leader update error: %v", err)
				}

				err = peer.Close()
				if err != nil {
					log.Fatal("Closing connection error: ", err)
				}

				return
			}

			peer, err = utils.DialTimeout("tcp", currentNode.List.GetNode(currentNode.List.Nodes[pass].Id).Address,
				5*time.Second)
			info.SkipCount = i
			if err != nil {
				continue
			} else {
				break
			}
		}
	}

	err = peer.Call("PeerServiceHandler.ElectionMessageCR", info, nil)
	if err != nil {
		log.Fatal("Election message forward failed: ", err)
	}

	err = peer.Close()
	if err != nil {
		log.Fatal("Closing connection error: ", err)
	}
}

func ChangAndRoberts(currNode utils.NodeINFO) {
	if len(currNode.List.GetAllNodes()) == 1 || currNode.Id == currNode.List.GetNode(currNode.Leader).Id {
		return
	}

	/* Performed only when new peer enter the system because it has Leader = -1 */
	if currNode.Id > currNode.Leader {
		fmt.Println("--- Start new election ---")

		/* Perform first iteration of algorithm */
		firstIteration(currNode)

		return
	}

	/* Attempt to ping leader process */
	peer, err := utils.DialTimeout("tcp", currNode.List.GetNode(currNode.Leader).Address, 5*time.Second)
	if err != nil {
		fmt.Println("--- Start new election ---")

		/* Perform first iteration of algorithm */
		firstIteration(currNode)

		return
	}

	defer func(peer *rpc.Client) {
		err = peer.Close()
		if err != nil {
			log.Fatal("Closing connection error: ", err)
		}
	}(peer)

	err = peer.Call("PeerServiceHandler.CheckLeaderStatus", currNode, nil)
	if err != nil {
		log.Printf("Ping to leader failed: %v\n", err)
	}
}

// firstIteration compute first iteration of algorithm
func firstIteration(currNode utils.NodeINFO) {
	mex := utils.Message{StartingMex: true}
	peer, err := utils.DialTimeout("tcp", currNode.List.GetNode(currNode.Id).Address, 5*time.Second)

	err = peer.Call("PeerServiceHandler.ElectionMessageCR", mex, nil)
	if err != nil {
		log.Printf("Leader update error: %v", err)
	}

	err = peer.Close()
	if err != nil {
		log.Fatal("Closing connection error: ", err)
	}
}
