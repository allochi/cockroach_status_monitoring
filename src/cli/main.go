package main

import (
	"encoding/json"
	"flag"
	"log"
	"regexp"

	models "../models"
	providers "../providers"
)

func main() {
	var entryAddress string
	flag.StringVar(&entryAddress, "address", "", "entry host http address in the format `address:port` or left empty for `localhost:8080")
	flag.Parse()

	if entryAddress != "" {
		if !isValidAddress(entryAddress) {
			log.Fatalf("host address must be in the format `address:port` or left empty for `localhost:8080`\n")
		}
	} else {
		entryAddress = "localhost:8080"
	}

	ch := models.ClusterHealth{
		EntryAddress:   entryAddress,
		StatusProvider: providers.HTTPStatusProvider{Address: entryAddress},
	}

	ch.Update()

	for _, node := range ch.Nodes {
		content, _ := json.Marshal(node)
		log.Println(string(content))
	}
}

func isValidAddress(address string) bool {
	match, _ := regexp.MatchString(`([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\:[0-9]{5})|(localhost\:[0-9]{5})|([a-z]+\.[a-z]+\.[a-z]+\:[0-9]{5})`, address)
	return match
}
