package algorithm

import (
	"SDCCproject/utils"
	"fmt"
	"log"
	"net/rpc"
)

func election(currNode utils.NodeINFO, actualLeader utils.Node) {
	for _, node := range currNode.List.GetAllNodes() {
		if node.Id > currNode.Id {
			peer, err := rpc.Dial("tcp", node.Address)
			if err != nil {
				continue
			}

			var repOK string
			err = peer.Call("NodeListUpdate.ElectionMessage", currNode, &repOK)
			if err != nil {
				log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
			}

			fmt.Println("OK message replied")

			err = peer.Close()
			if err != nil {
				log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
			}

			return
		}
	}

	fmt.Printf("You're the new leader!")
	for _, node := range currNode.List.GetAllNodes() {
		peer, err := rpc.Dial("tcp", node.Address)
		if err != nil {
			continue
		}

		leaderINFO := utils.LeaderStatus{NewLeaderID: currNode.Id, OldLeaderID: actualLeader.Id}

		err = peer.Call("NodeListUpdate.NewLeader", leaderINFO, nil)
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
	for _, node := range currNode.List.GetAllNodes() {
		if node.Leader == true {
			peer, err := rpc.Dial("tcp", node.Address)

			if err != nil {
				fmt.Println("--- Start new election ---")
				election(currNode, node)
				fmt.Println("--- Leader election terminated ---")
				break
			}

			err = peer.Call("NodeListUpdate.CheckLeaderStatus", node, nil)
			if err != nil {
				log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
			}

			err = peer.Close()
			if err != nil {
				log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
			}

			return
		}
	}
}
