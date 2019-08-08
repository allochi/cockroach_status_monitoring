package main

import (
	"io/ioutil"
	"testing"
)

func TestParseNodes(t *testing.T) {
	content, err := ioutil.ReadFile("mocks/status_nodes_response.json")
	nodes, err := ParseNodes(content)

	if err != nil {
		t.Error("failed to parse well formatted response")
	}

	if len(nodes) != 4 {
		t.Errorf("failed to parse well formatted response, expected 4 nodes for %d", len(nodes))
	}

	// test for default http address
	if nodes[0].HTTPAddress != "localhost:8080" {
		t.Errorf("failed to retrieve http address, expected 'localhost:8080' got '%s'", nodes[0].HTTPAddress)
	}
}
