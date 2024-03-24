package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"strconv"
)

func GetAddress() (string, error) {
	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	var ip string
	for _, addr := range interfaceAddr {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ip = ipNet.IP.String()
			break
		}
	}

	if ip == "" {
		return "", fmt.Errorf("local IP address not found")
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", fmt.Errorf("available port not found: %v", err)
	}

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {

		}
	}(listener)

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
