package algorithm

import (
	"SDCCproject/utils"
	"fmt"
	"log"
	"net/rpc"
)

func WinnerMessage(currentNode utils.NodeINFO, leader int) {
	info := utils.Message{
		SkipCount: 1,
		MexID:     leader,
		CurrNode:  currentNode,
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

	err = peer.Call("NodeListUpdate.NewLeaderCR", info, nil)
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
		info = utils.Message{SkipCount: 1, MexID: currentNode.Id, CurrNode: currentNode}
	} else {
		info = utils.Message{SkipCount: 1, MexID: mexReply, CurrNode: currentNode}
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

	err = peer.Close()
	if err != nil {
		log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
	}
}

func ChangAndRobert(currNode utils.NodeINFO) {
	peer, err := rpc.Dial("tcp", currNode.List.GetNode(currNode.Leader).Address)
	if err != nil {
		fmt.Println("--- Start new election ---")
		ElectionChangRobert(currNode, 0)
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

	return
}
