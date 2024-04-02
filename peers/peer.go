package main

import (
	"SDCCproject/algorithm"
	"SDCCproject/utils"
	"time"

	"fmt"
	"log"
	"net"
	"net/rpc"
)

// Global variable to hold peer status in DS
var currentNode utils.NodeINFO

func chooseAlgorithm() {
	config, err := utils.ReadConfig("/app/config.json")
	if err != nil {
		log.Fatal("Configuration file reading error: ", err)
	}

	for {
		/* Choose algorithm */
		if config.Algorithm.Bully == "true" {
			algorithm.Bully(currentNode)
		} else if config.Algorithm.Chang == "true" {
			algorithm.ChangAndRoberts(currentNode)
		} else {
			fmt.Println("Configuration file format wrong!")
		}

		/* Delay and print actual leader */
		fmt.Printf("Leader: %d\n", currentNode.Leader)
		time.Sleep(time.Second) // To test Chang & Robert increase the delay value
	}
}

func main() {
	config, err := utils.ReadConfig("/app/config.json")

	/* Init peer's service */
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

	/* Register peer's service on Service Registry */
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

	fmt.Printf("Your ID: %d, your address: %s, Nodes in system: %s\n", currentNode.Id, currentNode.Address, currentNode.List.Nodes)

	// go utils.StopNode(currentNode) // To simulate the auto peer crash uncomment this line
	go chooseAlgorithm()

	/* Listen for RPC */
	for {
		conn, _ := list.Accept()
		go peer.ServeConn(conn)
	}
}
