package main

import (
	"SDCCproject/utils"
	"bufio"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
)

func getAddress() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	var ip string
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip = ipNet.IP.String()
			break
		}
	}

	if ip == "" {
		return "", fmt.Errorf("indirizzo IP locale non trovato")
	}

	// Ottieni una porta disponibile
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", fmt.Errorf("impossibile ottenere una porta disponibile: %v", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

	// Ottieni l'indirizzo e la porta dell'ascoltatore
	addr := listener.Addr().(*net.TCPAddr)

	return ip + ":" + strconv.Itoa(addr.Port), nil
}

func keyboardInput() string {
	var scanner *bufio.Scanner

	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Fatal("Errore nell'acquisizione dell'input: ", err)
	}
	return scanner.Text()
}

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

func main() {
	address, err := getAddress()
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

		/* Bully algorithm */
		go bully()
	}
}

func bully() {
	fmt.Print("Start bully algorithm? Reply with Y to continue: ")

	if keyboardInput() == "Y" {
		for _, node := range currentNode.List.GetAllNodes() {
			if node.Leader == true {
				peer, err := rpc.Dial("tcp", node.Address)

				// TODO: qui capisco che il nodo leader non risponde più quindi devo indire la nuova elezione e gestire il crash
				if err != nil {
					log.Printf("Errore di connessione: %v", err)
				}

				err = peer.Call("NodeListUpdate.CheckLeaderStatus", node, nil)
				if err != nil {
					log.Printf("Errore durante l'aggiornamento del nodo: %v", err)
				}

				err = peer.Close()
				if err != nil {
					log.Fatalf("Errore durante l'aggiornamento del nodo: %v\n", err)
				}
			}
		}
	} else {
		fmt.Println("Invalid input!!!")
	}
}
