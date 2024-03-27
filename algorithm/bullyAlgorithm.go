package algorithm

import (
	"SDCCproject/utils"
	"fmt"
	"log"
	"net/rpc"
)

func ElectionBully(currNode utils.NodeINFO) {
	for _, node := range currNode.List.GetAllNodes() {
		if node.Id > currNode.Id {
			peer, err := rpc.Dial("tcp", node.Address)
			if err != nil {
				continue
			}

			var repOK string
			err = peer.Call("PeerServiceHandler.ElectionMessageBULLY", currNode, &repOK)
			if err != nil {
				log.Fatal("Election message forward failed: ", err)
			}

			err = peer.Close()
			if err != nil {
				log.Fatal("Closing connection error: ", err)
			}

			if repOK != "" {
				return
			}
		}
	}

	for _, node := range currNode.List.GetAllNodes() {
		peer, err := rpc.Dial("tcp", node.Address)
		if err != nil {
			continue
		}

		err = peer.Call("PeerServiceHandler.NewLeader", currNode.List.GetNode(currNode.Id), nil)
		if err != nil {
			log.Fatal("Leader update error: ", err)
		}

		err = peer.Close()
		if err != nil {
			log.Fatal("Closing connection error: ", err)
		}
	}
}

func Bully(currNode utils.NodeINFO) {
	if len(currNode.List.GetAllNodes()) == 1 || currNode.Id == currNode.List.GetNode(currNode.Leader).Leader {
		return
	}

	if currNode.Id > currNode.Leader {
		fmt.Println("--- Start new election ---")
		ElectionBully(currNode)
		return
	}

	/* Attempt to ping leader process */
	peer, err := rpc.Dial("tcp", currNode.List.GetNode(currNode.Leader).Address)
	if err != nil {
		fmt.Println("--- Start new election ---")
		ElectionBully(currNode)
		return
	}

	defer func(peer *rpc.Client) {
		err = peer.Close()
		if err != nil {
			log.Fatal("Closing connection error: ", err)
		}
	}(peer)

	/* Call CheckLeaderStatus on leader */
	err = peer.Call("PeerServiceHandler.CheckLeaderStatus", currNode.List.GetNode(currNode.Leader), nil)
	if err != nil {
		log.Fatal("Ping to leader failed: ", err)
	}
}
