package main

import (
	"SDCCproject/algorithm"
	"SDCCproject/utils"

	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
)

// Global variable to hold peers in DS
var currentNode utils.NodeINFO

// PeerServiceHandler is the interface which exposes the RPC method ManageNode
type PeerServiceHandler struct{}

func (PeerServiceHandler) UpdateList(node utils.Node, _ *utils.NodeINFO) error {
	currentNode.List.GetAllNodes()
	currentNode.List.AddNode(node)

	fmt.Printf("New peer in system, list updated\n")

	return nil
}

func (PeerServiceHandler) CheckLeaderStatus(_ utils.Node, _ *utils.NodeINFO) error {
	fmt.Printf("Hi i'm the leader peer with address: %s, and id: %d. I'm still active!\n", currentNode.Address, currentNode.Id)

	return nil
}

func (PeerServiceHandler) ElectionMessageBULLY(nodeCaller utils.NodeINFO, rep *string) error {
	fmt.Printf("Election message from peer with id: %d\n", nodeCaller.Id)
	*rep = "OK"

	go algorithm.ElectionBully(currentNode)

	return nil
}

func (PeerServiceHandler) NewLeader(leaderNode utils.NodeINFO, _ *utils.NodeINFO) error {
	currentNode.Leader = leaderNode.Id

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

	for _, node := range currentNode.List.GetAllNodes() {
		currentNode.List.UpdateNode(node, mex.MexID)
	}

	go algorithm.WinnerMessage(currentNode.List.GetNode(nextNode), mex.MexID)

	return nil
}

func (PeerServiceHandler) ElectionMessageCR(mex utils.Message, _ *int) error {
	currID := (mex.CurrNode.Id + mex.SkipCount) % len(currentNode.List.Nodes)

	if mex.MexID > currID {
		go algorithm.ElectionChangRobert(currentNode.List.GetNode(currID), mex.MexID)
	} else if mex.MexID < currID {
		mex.MexID = currID
		go algorithm.ElectionChangRobert(currentNode.List.GetNode(currID), mex.MexID)
	} else if mex.MexID == currID {
		info := utils.Message{
			SkipCount: 1,
			MexID:     mex.MexID,
			CurrNode:  currentNode.List.GetNode(currID),
		}

		currentNode.Leader = mex.MexID

		peer, err := rpc.Dial("tcp", currentNode.List.GetNode((currentNode.Id+1)%len(currentNode.List.Nodes)).Address)
		if err != nil {
			skip := (currentNode.Id + 1) % len(currentNode.List.Nodes)
			i := 1
			for {
				i++
				pass := (currentNode.Id + i) % len(currentNode.List.Nodes)
				if pass == skip-1 {
					return nil
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

func chooseAlgorithm() {
	for {
		/* Choose algorithm */
		//algorithm.Bully(currentNode)
		algorithm.ChangAndRobert(currentNode)
	}
}

func stopNode() {
	minNum := 0
	maxNum := 100000

	for {
		randNum := utils.Random(minNum, maxNum)
		if currentNode.Id == randNum {
			os.Exit(1)
		}
	}
}

func printLeader() {
	for {
		fmt.Printf("Leader: %d\n", currentNode.Leader)
	}
}

func main() {
	address, err := utils.GetAddress()
	if err != nil {
		fmt.Println("IP retrieving error: ", err)
	}

	/* Init peer's service */
	peerService := new(PeerServiceHandler)

	peer := rpc.NewServer()
	err = peer.Register(peerService)
	if err != nil {
		log.Fatal("Wrong service format: ", err)
	}

	list, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	/* Register peer's service on Service Registry */
	// Reading configuration file
	config, err := utils.ReadConfig("/Users/andreaandreoli/Desktop/projectSDCC/config.json")
	if err != nil {
		log.Fatal("Configuration file reading error: ", err)
	}

	serviceAddress := config.ServiceRegistry.Address + config.ServiceRegistry.Port
	server, err := rpc.Dial("tcp", serviceAddress)
	if err != nil {
		log.Fatal("Connection error SR: ", err)
	}

	defer func(server *rpc.Client) {
		err := server.Close()
		if err != nil {
			log.Fatal("Closing connection error: ", err)
		}
	}(server)

	/* Entering distributed system with RPC call to the service registry */
	peerRep := utils.NodeINFO{}
	peerInfo := utils.PeerAddr{PeerAddress: address}

	err = server.Call("NodeHandler.ManageNode", peerInfo, &peerRep)
	if err != nil {
		log.Fatal("Entering System error: ", err)
	}

	fmt.Printf("Peer INFO [ID: %d, Leader: %t, ListOfNode: %s, Address: %s]\n", peerRep.Id, peerRep.Leader, peerRep.List, address)
	currentNode = peerRep
	currentNode.Address = address

	go stopNode()
	go chooseAlgorithm()
	go printLeader() // active to test new leader

	/* Listen for RPC */
	for {
		conn, _ := list.Accept()
		go peer.ServeConn(conn)
	}
}
