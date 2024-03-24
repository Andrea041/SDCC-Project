package main

import (
	"SDCCproject/utils"

	"log"
	"net"
	"net/rpc"
)

// NodeHandler is the interface which exposes the RPC method ManageNode
type NodeHandler struct{}

// Global variable to hold DockerfilePeer in DS
var nodeList utils.NodeList

// ManageNode is the discovery service for each new DockerfilePeer in DS
func (NodeHandler) ManageNode(peerAddress utils.PeerAddr, nodeInfo *utils.NodeINFO) error {
	nodes := nodeList.GetAllNodes()

	// Find max ID
	maxID := -1
	for _, node := range nodes {
		if node.Id > maxID {
			maxID = node.Id
		}
	}
	newID := maxID + 1

	var newNode utils.Node
	newNode = utils.Node{Id: newID, Address: peerAddress.PeerAddress, Leader: 0}

	// Add new node to local list
	nodeList.AddNode(newNode)

	// Send DockerfilePeer's info
	nodeInfo.Id = newID
	nodeInfo.List.Nodes = nodeList.GetAllNodes()
	nodeInfo.Leader = newNode.Leader

	for _, node := range nodes {
		peer, err := rpc.Dial("tcp", node.Address)
		if err != nil {
			continue
		}

		err = peer.Call("PeerServiceHandler.UpdateList", newNode, nil)
		if err != nil {
			log.Fatal("List update failed: ", err)
		}

		err = peer.Close()
		if err != nil {
			log.Fatal("Closing connection error: ", err)
		}
	}

	return nil
}

func main() {
	serviceRegistry := new(NodeHandler)

	server := rpc.NewServer()
	err := server.Register(serviceRegistry)
	if err != nil {
		log.Fatal("Wrong service format: ", err)
	}

	// TODO: Scrivere nel report finale che tra le ipotesi dell'algoritmo si ha comunicazione affidabile quindi uso TCP come protocollo di comunicazione
	config, err := utils.ReadConfig("/Users/andreaandreoli/Desktop/projectSDCC/config.json")
	if err != nil {
		log.Fatal("Configuration file reading error: ", err)
	}

	serviceAddress := config.ServiceRegistry.Address + config.ServiceRegistry.Port
	list, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	log.Printf("Service registry listening on port %s", config.ServiceRegistry.Port)

	for {
		conn, _ := list.Accept()
		go server.ServeConn(conn)
	}
}
