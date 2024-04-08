package algorithm

import (
	"SDCCproject/utils"

	"fmt"
	"log"
	"net/rpc"
	"time"
)

func ElectionBully(currNode utils.NodeINFO) {
	for _, node := range currNode.List.GetAllNodes() {
		if node.Id > currNode.Id {
			peer, err := utils.DialTimeout("tcp", node.Address, 5*time.Second)
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

			/* Test if the successor node reply */
			if repOK != "" {
				/* Don't call other processes: one bigger ID is ac */
				return
			}
		}
	}

	for _, node := range currNode.List.GetAllNodes() {
		peer, err := utils.DialTimeout("tcp", node.Address, 5*time.Second)
		if err != nil {
			continue
		}

		err = peer.Call("PeerServiceHandler.NewLeaderBULLY", currNode.List.GetNode(currNode.Id), nil)
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
	if len(currNode.List.GetAllNodes()) == 1 || currNode.Id == currNode.List.GetNode(currNode.Leader).Id {
		return
	}

	/* Performed only when new peer enter the system because it has leader = -1 */
	if currNode.Id > currNode.Leader {
		fmt.Println("--- Start new election ---")
		ElectionBully(currNode)
		return
	}

	/* Attempt to ping leader process */
	peer, err := utils.DialTimeout("tcp", currNode.List.GetNode(currNode.Leader).Address, 5*time.Second)
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
	err = peer.Call("PeerServiceHandler.CheckLeaderStatus", currNode, nil)
	if err != nil {
		log.Fatal("Ping to leader failed: ", err)
	}
}
