package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// client, err := http.Get("http://localhost:8080/_status/nodes")
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// body, err := ioutil.ReadAll(client.Body)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// // var state map[string]interface{}
	// var state struct {
	// 	Nodes []struct {
	// 		Desc struct {
	// 			NodeID  int
	// 			Address struct {
	// 				AddressField string
	// 			}
	// 		}
	// 		StartedAt     string
	// 		UpdatedAt     string
	// 		StoreStatuses []struct {
	// 			Metrics struct {
	// 				Capacity          int64
	// 				CapacityAvaliable int64 `json:"capacity.available"`
	// 			}
	// 		}
	// 		Args []string
	// 	}
	// }
	// err = json.Unmarshal(body, &state)
	// if err != nil {
	// 	log.Fatalln("couldn't read state: ", err)
	// }

	// fmt.Println(state.Nodes)
	entryAddress := "localhost:8080"
	client, err := http.Get("http://" + entryAddress + "/_status/nodes")
	if err != nil {
		// 	return []Node{}, err
	}

	body, err := ioutil.ReadAll(client.Body)
	// if err != nil {
	// 	return []Node{}, err
	// }

	nodes, _ := ParseNodes(body)
	fmt.Println(nodes)
}

type Node struct {
	NodeID            int
	Address           string
	HTTPAddress       string
	StartedAt         string
	UpdatedAt         string
	Capacity          int64
	CapacityAvaliable int64
}

func ParseNodes(content []byte) ([]Node, error) {
	var state struct {
		Nodes []struct {
			Desc struct {
				NodeID  int
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
		return []Node{}, err
	}

	var nodes []Node
	for _, n := range state.Nodes {
		var node Node
		node.NodeID = n.Desc.NodeID
		node.Address = n.Desc.Address.AddressField
		node.StartedAt = n.StartedAt
		node.UpdatedAt = n.UpdatedAt
		node.Capacity = n.StoreStatuses[0].Metrics.Capacity
		node.CapacityAvaliable = n.StoreStatuses[0].Metrics.CapacityAvaliable

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
