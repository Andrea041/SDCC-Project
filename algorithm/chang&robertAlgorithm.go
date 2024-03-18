package algorithm

import (
	"SDCCproject/utils"
	"fmt"
	"log"
	"net/rpc"
)

var actualLeader utils.Node

// TODO: Sistemare qui
func WinnerMessage(currentNode utils.NodeINFO) {
	if currentNode.Participant == false {
		return
	} else {
		currentNode.Participant = false
	}

	peer, err := rpc.Dial("tcp", currentNode.List.GetNode((currentNode.Id+1)%len(currentNode.List.Nodes)).Address)
	if err != nil {
		log.Printf("Errore di connessione: %v", err)
	}

	leaderINFO := utils.LeaderStatus{NewLeaderID: currentNode.Id, OldLeaderID: actualLeader.Id}

	err = peer.Call("NodeListUpdate.NewLeaderCR", leaderINFO, currentNode)
	if err != nil {
		log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
	}

	err = peer.Close()
	if err != nil {
		log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
	}
}

func ElectionChangRobert(currentNode utils.NodeINFO, mexReply int) {
	var info utils.Message

	if mexReply == 0 {
		info = utils.Message{SkipCount: 1, CurrNode: currentNode, MexID: currentNode.Id}
	} else {
		info = utils.Message{SkipCount: 1, CurrNode: currentNode, MexID: mexReply}
	}

	peer, err := rpc.Dial("tcp", currentNode.List.GetNode((currentNode.Id+1)%len(currentNode.List.Nodes)).Address)
	if err != nil {
		skip := (currentNode.Id + 1) % len(currentNode.List.Nodes)
		i := skip
		for {
			i++
			pass := i % len(currentNode.List.Nodes)
			if pass == skip-1 {
				return
			}

			peer, err = rpc.Dial("tcp", currentNode.List.GetNode((currentNode.Id+i)%len(currentNode.List.Nodes)).Address)
			info.SkipCount = i
			if err != nil {
				continue
			} else {
				break
			}
		}
	}

	err = peer.Call("NodeListUpdate.ElectionMessageCR", info, nil)
	if err != nil {
		log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
	}

	currentNode.Participant = true

	err = peer.Close()
	if err != nil {
		log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
	}
}

func ChangAndRobert(currNode utils.NodeINFO) {
	for _, node := range currNode.List.GetAllNodes() {
		if node.Leader == true {
			actualLeader = node
			peer, err := rpc.Dial("tcp", node.Address)

			if err != nil {
				fmt.Println("--- Start new election ---")
				ElectionChangRobert(currNode, 0)
				fmt.Println("--- Leader election terminated ---")
				fmt.Println("breaking")
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
	fmt.Println("exiting")
}
