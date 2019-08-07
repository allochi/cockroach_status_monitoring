package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	client, err := http.Get("http://localhost:8080/_status/vars")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(client.Body)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(isLowMemory(string(body)))
}

func isLowMemory(content string) bool {
	var capacity float64
	var capacityAvailable float64
	scanner := bufio.NewScanner(strings.NewReader(content))
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
