package main

import (
	"SDCCproject/algorithm"
	"SDCCproject/utils"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"time"
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
	fmt.Printf("Hi i'm the leader with address: %s, and id: %d. I'm still active!\n", currentNode.Address, currentNode.Id)

	return nil
}

func (PeerServiceHandler) ElectionMessageBULLY(nodeCaller utils.NodeINFO, rep *string) error {
	fmt.Printf("Election message from DockerfilePeer with id: %d\n", nodeCaller.Id)
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
		algorithm.Bully(currentNode)
		//algorithm.ChangAndRobert(currentNode)
		time.Sleep(time.Second)
		fmt.Printf("Leader: %d\n", currentNode.Leader)
	}
}

func stopNode() {
	minNum := 0
	maxNum := 100000000

	for {
		randNum := utils.Random(minNum, maxNum)
		if currentNode.Id == randNum {
			os.Exit(1)
		}
	}
}

func main() {
	config, err := utils.ReadConfig("/app/config.json")
	/* Init DockerfilePeer's service */
	peerService := new(PeerServiceHandler)

	peer := rpc.NewServer()
	err = peer.Register(peerService)
	if err != nil {
		log.Fatal("Wrong service format: ", err)
	}

	peerAddress := config.Peer.Address + config.Peer.Port
	list, err := net.Listen("tcp", peerAddress)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	address := list.Addr().String()

	/* Register DockerfilePeer's service on Service Registry */
	//config, err := utils.ReadConfig("/Users/andreaandreoli/Desktop/projectSDCC/config.json")
	if err != nil {
		log.Fatal("Configuration file reading error: ", err)
	}

	serviceAddress := config.ServiceRegistry.Address + config.ServiceRegistry.Port
	server, err := rpc.Dial("tcp", serviceAddress)
	if err != nil {
		log.Fatal("Connection error to service registry: ", err)
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

	currentNode = peerRep
	currentNode.Address = address

	fmt.Printf("Your ID: %d, leader ID: %d, your address: %s, Nodes in system: %s\n", currentNode.Id, currentNode.Leader, currentNode.Address, currentNode.List.Nodes)

	// go stopNode()
	go chooseAlgorithm()

	/* Listen for RPC */
	for {
		conn, _ := list.Accept()
		go peer.ServeConn(conn)
	}
}
