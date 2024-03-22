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
			err = peer.Call("NodeListUpdate.ElectionMessageBULLY", currNode, &repOK)
			if err != nil {
				log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
			}

			if repOK != "" {
				fmt.Println("OK message replied")

				err = peer.Close()
				if err != nil {
					log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
				}

				return
			}
		}
	}

	for _, node := range currNode.List.GetAllNodes() {
		peer, err := rpc.Dial("tcp", node.Address)
		if err != nil {
			continue
		}

		err = peer.Call("NodeListUpdate.NewLeader", currNode.List.GetNode(currNode.Id), nil)
		if err != nil {
			log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
		}

		err = peer.Close()
		if err != nil {
			log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
		}
	}
}

func Bully(currNode utils.NodeINFO) {
	if len(currNode.List.GetAllNodes()) == 1 {
		return
	}

	if currNode.Id == currNode.List.GetNode(currNode.Leader).Leader {
		return
	}

	if currNode.Id > currNode.Leader {
		ElectionBully(currNode)
	}

	/* Ping leader process */
	peer, err := rpc.Dial("tcp", currNode.List.GetNode(currNode.Leader).Address)
	if err != nil {
		fmt.Println("--- Start new election ---")
		ElectionBully(currNode)
		return
	}

	err = peer.Call("NodeListUpdate.CheckLeaderStatus", currNode.List.GetNode(currNode.Leader), nil)
	if err != nil {
		log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
	}

	err = peer.Close()
	if err != nil {
		log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
	}
	fmt.Println("Connection closed")
}
