package utils

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"
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
