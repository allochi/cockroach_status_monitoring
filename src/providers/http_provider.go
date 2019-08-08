package main

import (
	"bufio"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// IsNodeLowInMemory checks if a node os low in memory
func IsNodeLowInMemory(address string) (bool, error) {
	client, err := http.Get(address)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(client.Body)
	if err != nil {
		return false, err
	}
	logs := string(body)

	return isLowInMemory(logs), nil
}

func isLowInMemory(logs string) bool {
	var capacity float64
	var capacityAvailable float64
	scanner := bufio.NewScanner(strings.NewReader(logs))
	for scanner.Scan() {
		text := scanner.Text()

		// find capacity
		if strings.HasPrefix(text, "capacity{") {
			capacity = getMemorySizeFromLog(text)
		}

		// find capacity_available
		if strings.HasPrefix(text, "capacity_available{") {
			capacityAvailable = getMemorySizeFromLog(text)
		}
	}

	return capacityAvailable == 0 || (capacity/capacityAvailable < 0.15)
}

func getMemorySizeFromLog(line string) float64 {
	parts := strings.Split(line, " ")
	var size float64
	var err error
	if len(parts) == 2 {
		size, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			size = 0
		}
	}
	return size
}
