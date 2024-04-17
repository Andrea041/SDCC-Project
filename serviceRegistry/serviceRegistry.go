package main

import (
	"SDCCproject/utils"
	"log"
	"net"
	"net/rpc"
)

// Global variable to hold Peer in DS
var nodeList utils.NodeList

func main() {
	serviceRegistry := new(NodeHandler)

	server := rpc.NewServer()
	err := server.Register(serviceRegistry)
	if err != nil {
		log.Fatal("Wrong service format: ", err)
	}

	config, err := utils.ReadConfigJSON("/app/config.json")
	if err != nil {
		log.Fatal("Configuration file reading error: ", err)
	}

	serviceAddress := config.ServiceRegistry.Address + config.ServiceRegistry.Port
	list, err := net.Listen("tcp", serviceAddress)
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	log.Printf("Service registry listening on %s%s", config.ServiceRegistry.Address, config.ServiceRegistry.Port)

	for {
		conn, _ := list.Accept()
		go server.ServeConn(conn)
	}
}
