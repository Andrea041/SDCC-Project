package utils

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net"
	"net/rpc"
	"os"
	"time"
)

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

func DialTimeout(network string, address string, timeout time.Duration) (*rpc.Client, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	if err != nil {
		return nil, err
	}

	return rpc.NewClient(conn), nil
}

func StopNode(currentNode NodeINFO) {
	minNum := 0
	maxNum := 100000000

	for {
		randNum := Random(minNum, maxNum)
		if currentNode.Id == randNum {
			os.Exit(1)
		}
	}
}
