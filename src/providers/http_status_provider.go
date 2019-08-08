package providers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	models "../models"
)

// HTTPStatusProvider HTTP status provider
type HTTPStatusProvider struct {
	Address string
}

// Call calls a node and return the status of all nodes
func (sp HTTPStatusProvider) Call() ([]models.Node, error) {
	if sp.Address == "" {
		return []models.Node{}, fmt.Errorf("http host address is required")
	}

	client, err := http.Get("http://" + sp.Address + "/_status/nodes")
	if err != nil {
		return []models.Node{}, err
	}

	body, err := ioutil.ReadAll(client.Body)
	if err != nil {
		return []models.Node{}, err
	}

	nodes, err := sp.parseNodes(body)
	if err != nil {
		return []models.Node{}, err
	}

	// check each node is alive
	for i := range nodes {
		client, err := http.Get("http://" + nodes[i].HTTPAddress + "/health")
		nodes[i].IsLive = err == nil && client.StatusCode == http.StatusOK
	}

	// check each node is available
	for i := range nodes {
		if nodes[i].IsLive {
			client, err := http.Get("http://" + nodes[i].HTTPAddress + "/health?ready=1")
			nodes[i].IsAvailable = err == nil && client.StatusCode == http.StatusOK
		}
	}

	return nodes, nil
}

func (sp HTTPStatusProvider) parseNodes(content []byte) ([]models.Node, error) {
	var state struct {
		Nodes []struct {
			Desc struct {
				NodeID  int64
				Address struct {
					AddressField string
				}
			}
			StartedAt     string
			UpdatedAt     string
			StoreStatuses []struct {
				Metrics struct {
					Capacity          int64
					CapacityAvaliable int64 `json:"capacity.available"`
				}
			}
			Args []string
		}
	}
	err := json.Unmarshal(content, &state)
	if err != nil {
		return []models.Node{}, err
	}

	var nodes []models.Node
	for _, n := range state.Nodes {
		var node models.Node
		node.ID = n.Desc.NodeID
		node.Address = n.Desc.Address.AddressField
		node.StartedAt = n.StartedAt
		node.UpdatedAt = n.UpdatedAt
		node.Capacity = n.StoreStatuses[0].Metrics.Capacity
		node.CapacityAvaliable = n.StoreStatuses[0].Metrics.CapacityAvaliable
		node.IsLowInMemory = node.Capacity == 0 || node.CapacityAvaliable == 0 || float64(node.Capacity)/float64(node.CapacityAvaliable) < 0.15

		// http address
		for _, arg := range n.Args {
			if strings.HasPrefix(arg, "--http-addr") {
				parts := strings.Split(arg, "=")
				node.HTTPAddress = parts[1]
			}
		}

		// set default http address
		if node.HTTPAddress == "" {
			parts := strings.Split(node.Address, ":")
			node.HTTPAddress = parts[0] + ":8080"
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}
