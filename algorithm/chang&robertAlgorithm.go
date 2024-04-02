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

	peer, err := utils.DialTimeout("tcp", currentNode.List.GetNode((currentNode.Id+1)%len(currentNode.List.Nodes)).Address, 5*time.Second)
	if err != nil {
		skip := (currentNode.Id + 1) % len(currentNode.List.Nodes)
		i := 0
		for {
			i++
			pass := (currentNode.Id + i) % len(currentNode.List.Nodes)
			if pass == skip-1 {
				return
			}

			peer, err = utils.DialTimeout("tcp", currentNode.List.GetNode((currentNode.Id+i)%len(currentNode.List.Nodes)).Address, 5*time.Second)
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

func ElectionChangRobert(currentNode utils.NodeINFO, mexReply int) {
	var info utils.Message

	info = utils.Message{SkipCount: 1, MexID: mexReply, CurrNode: currentNode}

	peer, err := utils.DialTimeout("tcp", currentNode.List.GetNode((currentNode.Id+1)%len(currentNode.List.Nodes)).Address, 5*time.Second)
	if err != nil {
		skip := (currentNode.Id + 1) % len(currentNode.List.Nodes)
		i := 0
		for {
			i++
			pass := (currentNode.Id + i) % len(currentNode.List.Nodes)
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

			peer, err = utils.DialTimeout("tcp", currentNode.List.GetNode((currentNode.Id+i)%len(currentNode.List.Nodes)).Address, 5*time.Second)
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

func ChangAndRobert(currNode utils.NodeINFO) {
	if len(currNode.List.GetAllNodes()) == 1 || currNode.Id == currNode.List.GetNode(currNode.Leader).Leader {
		if currNode.Leader == -1 {
			currNode.Leader = currNode.Id
		}
		return
	}

	if currNode.Id > currNode.Leader {
		fmt.Println("--- Start new election ---")
		ElectionChangRobert(currNode, currNode.Id)
		return
	}

	peer, err := utils.DialTimeout("tcp", currNode.List.GetNode(currNode.Leader).Address, 5*time.Second)
	if err != nil {
		fmt.Println("--- Start new election ---")
		ElectionChangRobert(currNode, currNode.Id)
		return
	}

	defer func(peer *rpc.Client) {
		err = peer.Close()
		if err != nil {
			log.Fatal("Closing connection error: ", err)
		}
	}(peer)

	err = peer.Call("PeerServiceHandler.CheckLeaderStatus", currNode.List.GetNode(currNode.Leader), nil)
	if err != nil {
		log.Printf("Ping to leader failed: %v\n", err)
	}
}
