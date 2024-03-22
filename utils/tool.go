package utils

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strconv"
)

func KeyboardInput() string {
	var scanner *bufio.Scanner

	scanner = bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		log.Fatal("Errore nell'acquisizione dell'input: ", err)
	}
	return scanner.Text()
}

func GetAddress() (string, error) {
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

func ReadConfig(file string) (Configuration, error) {
	var config Configuration

	tmp, err := os.ReadFile(file)
	if err != nil {
		return Configuration{}, err
	}

	err = json.Unmarshal(tmp, &config)
	if err != nil {
		return Configuration{}, err
	}
	return config, nil
}

func Random(min, max int) int {
	diff := big.NewInt(int64(max - min + 1))

	num, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return 0
	}

	return int(num.Int64()) + min
}
