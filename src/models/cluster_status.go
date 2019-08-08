package models

import (
	"log"
	"time"
)

// StatusProvider an interface for any provider that can return the status of the cluster nodes
type StatusProvider interface {
	Call() ([]Node, error)
}

// ClusterHealth a cache that holds the cluster health information
type ClusterHealth struct {
	EntryAddress   string // main host to call for information
	StatusProvider StatusProvider
	UpdatedAt      time.Time
	Nodes          []Node // cache status of nodes
}

// Update updates the cluster health from a provider
func (ch *ClusterHealth) Update() {
	// Call entry host
	httpNodes, err := ch.StatusProvider.Call()
	if err != nil {
		// Entry host failed find another host or exist if non available
		if len(ch.Nodes) == 0 {
			log.Fatalf("failed to connect to `%s` and there are no other alternative hosts registered", ch.EntryAddress)
		}

		var replaced bool
		for _, node := range ch.Nodes {
			ch.EntryAddress = node.Address
			httpNodes, err = ch.StatusProvider.Call()
			if err == nil {
				replaced = true
				break
			}
		}

		if !replaced {
			log.Fatalln("failed to connect to any of the registered hosts")
		}
	}

	ch.Nodes = httpNodes

	// TODO: Check live nodes with _status/vars

	ch.UpdatedAt = time.Now()
}
