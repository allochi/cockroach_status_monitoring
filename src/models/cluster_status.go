package models

import (
	"log"
	"time"
)

// StatusProvider an interface for any provider that can return the status of the cluster nodes
type StatusProvider interface {
	Call(string) ([]NodeStatus, error)
}

// NodeStatus type
// type NodeStatus struct {
// 	ID          string
// 	Address     string
// 	Build       string
// 	StartedAt   string
// 	UpdatedAt   string
// 	IsAvailable bool
// 	IsLive      bool
// }

// ClusterHealth a cache that holds the cluster health information
type ClusterHealth struct {
	EntryAddress string // main host to call for information
	Provider     StatusProvider
	UpdatedAt    time.Time
	Nodes        []NodeStatus // cache status of nodes
}

// Update updates the cluster health from a provider
func (ch *ClusterHealth) Update() {
	// Call entry host
	results, err := ch.Provider.Call(ch.EntryAddress)
	if err != nil {
		// Entry host failed find another host or exist if non available
		if len(ch.Nodes) == 0 {
			log.Fatalf("failed to connect to `%s` and there are no other alternative hosts registered", ch.EntryAddress)
		}

		for _, status := range ch.Nodes {
			results, err = ch.Provider.Call(ch.EntryAddress)
			if err == nil {
				ch.EntryAddress = status.Address
				break
			}
		}
		log.Fatalln("failed to connect to any of the registered hosts")
	}

	ch.Nodes = results
	ch.UpdatedAt = time.Now()
}
