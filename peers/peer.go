package main

import (
	"SDCCproject/algorithm"
	"SDCCproject/utils"

	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Global variable to hold peers in DS
var currentNode utils.NodeINFO

// NodeListUpdate is the interface which exposes the RPC method ManageNode
type NodeListUpdate struct{}

func (NodeListUpdate) UpdateList(node utils.Node, _ *utils.NodeINFO) error {
	currentNode.List.GetAllNodes()
	currentNode.List.AddNode(node)

	fmt.Printf("New peer in system, list updated\n")

	return nil
}

func (NodeListUpdate) CheckLeaderStatus(node utils.Node, _ *utils.NodeINFO) error {
	fmt.Printf("Hi i'm the peer with address: %s, with id: %d\n", node.Address, node.Id)

	return nil
}

func (NodeListUpdate) ElectionMessage(nodeCaller utils.NodeINFO, rep *string) error {
	fmt.Printf("Election message from peer with id: %d\n", nodeCaller.Id)
	*rep = "OK"

	go algorithm.Bully(currentNode)

	return nil
}

func (NodeListUpdate) NewLeader(leaderINFO utils.LeaderStatus, _ *utils.NodeINFO) error {
	for _, node := range currentNode.List.GetAllNodes() {
		if node.Id == leaderINFO.NewLeaderID {
			currentNode.List.UpdateNode(node, true)
		}

		if node.Id == leaderINFO.OldLeaderID {
			currentNode.List.UpdateNode(node, false)
		}
	}

	fmt.Printf("List updated: %s\n", currentNode.List.GetAllNodes())

	return nil
}

func main() {
	address, err := utils.GetAddress()
	if err != nil {
		fmt.Println("Errore durante il recupero dell'indirizzo IP locale:", err)
	}

	/* Init peer's service */
	peerService := new(NodeListUpdate)

	peer := rpc.NewServer()
	err = peer.Register(peerService)
	if err != nil {
		log.Fatal("Il formato del servizio è errato: ", err)
	}

	list, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Errore nell'instaurazione della connessione: ", err)
	}

	/* Register peer's service on Service Registry */
	// TODO: implementare successivamente con un file di configurazione il reperimento dell'indirizzo del service registry che per assuzione è noto, segnare nel report!!
	serviceReg := "localhost:" + "8888"
	server, err := rpc.Dial("tcp", serviceReg)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	defer func(server *rpc.Client) {
		err := server.Close()
		if err != nil {
			log.Fatal("Il formato del servizio è errato: ", err)
		}
	}(server)

	/* Entering distributed system with RPC call to the service registry */
	peerRep := utils.NodeINFO{}
	peerInfo := utils.PeerAddr{PeerAddress: address}

	err = server.Call("NodeHandler.ManageNode", peerInfo, &peerRep)
	if err != nil {
		log.Fatal("Entering System error: ", err)
	}

	fmt.Printf("Peer INFO [ID: %d, Leader: %t, Address: %s]\n", peerRep.Id, peerRep.Leader, peerRep.List)
	currentNode = peerRep

	/* Listen for RPC */
	for {
		conn, _ := list.Accept()
		go peer.ServeConn(conn)

		fmt.Println("--- Choose algorithm ---")
		fmt.Println("1 - Bully")
		fmt.Println("2 - Chang & Robert")
		go chooseAlgorithm()
	}
}

func chooseAlgorithm() {
	// Algorithm
	choose := utils.KeyboardInput()

	switch choose {
	case "1":
		go algorithm.Bully(currentNode)
	case "2":
		go algorithm.ChangAndRobert(currentNode)
	default:
		fmt.Println("Invalid input")
	}
}
