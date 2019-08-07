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
	var entryAddress = flag.String("address", "", "host address in the format `address:port` or left empty for `localhost:26257")
	flag.Parse()

	if *entryAddress != "" {
		if !isValidAddress(*entryAddress) {
			log.Fatalf("host address must be in the format `address:port` or left empty for `localhost:26257`\n")
		}
	} else {
		*entryAddress = "localhost:26257"
	}

	ch := models.ClusterHealth{
		EntryAddress: *entryAddress,
		Provider:     providers.CmdProvider{},
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
