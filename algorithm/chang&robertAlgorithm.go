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
		i := 0
		for {
			i++
			pass := (currentNode.Id + i) % len(currentNode.List.Nodes)
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

	info = utils.Message{SkipCount: 1, MexID: mexReply, CurrNode: currentNode}

	peer, err := rpc.Dial("tcp", currentNode.List.GetNode((currentNode.Id+1)%len(currentNode.List.Nodes)).Address)
	if err != nil {
		skip := (currentNode.Id + 1) % len(currentNode.List.Nodes)
		i := 0
		for {
			i++
			pass := (currentNode.Id + i) % len(currentNode.List.Nodes)
			if pass == skip-1 {
				info.MexID = currentNode.Id
				peer, err = rpc.Dial("tcp", currentNode.List.GetNode(currentNode.Id).Address)

				err = peer.Call("NodeListUpdate.NewLeaderCR", info, nil)
				if err != nil {
					log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
				}

				err = peer.Close()
				if err != nil {
					log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
				}

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
	if len(currNode.List.GetAllNodes()) == 1 {
		return
	}

	if currNode.Id == currNode.List.GetNode(currNode.Leader).Leader {
		return
	}

	if currNode.Id > currNode.Leader {
		ElectionChangRobert(currNode, currNode.Id)
	}

	peer, err := rpc.Dial("tcp", currNode.List.GetNode(currNode.Leader).Address)
	if err != nil {
		fmt.Println("--- Start new election ---")
		ElectionChangRobert(currNode, currNode.Id)
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
