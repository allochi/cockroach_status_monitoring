package main

import (
	"context"
	"flag"
	"log"
	"net"
	"regexp"
	"time"

	"google.golang.org/grpc"

	models "../models"
	providers "../providers"
)

var ch models.ClusterHealth

// StatusService gRPC service
type StatusService struct{}

// GetStatus gRPC implementation
func (s StatusService) GetStatus(ctx context.Context, void *models.Void) (*models.MonitoringResponse, error) {
	var nodes []*models.NodeStatus
	for _, node := range ch.Nodes {
		nodes = append(nodes, &node)
	}
	response := &models.MonitoringResponse{Nodes: nodes}
	return response, nil
}

func main() {

	var entryAddress = flag.String("address", "", "host address in the format `address:port` or left empty for `localhost:26257")
	var duration = flag.Int("duration", 5, "duration between each status update")
	flag.Parse()

	if *entryAddress != "" {
		if !isValidAddress(*entryAddress) {
			log.Fatalf("host address must be in the format `address:port` or left empty for `localhost:26257`\n")
		}
	} else {
		*entryAddress = "localhost:26257"
	}

	ch = models.ClusterHealth{
		EntryAddress: *entryAddress,
		Provider:     providers.CmdProvider{},
	}

	// Update cluster status periodically
	ticker := time.NewTicker(time.Duration(*duration) * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			ch.Update()
			// fmt.Printf("%#v\n", ch)
		}
	}()

	srv := grpc.NewServer()

	var statusService StatusService
	models.RegisterMonitoringServiceServer(srv, statusService)

	ln, err := net.Listen("tcp", ":8899")
	if err != nil {
		log.Fatalln("Couldn't listen on localhost:8899")
	}
	log.Fatal(srv.Serve(ln))

}

func isValidAddress(address string) bool {
	match, _ := regexp.MatchString(`([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\:[0-9]{5})|(localhost\:[0-9]{5})|([a-z]+\.[a-z]+\.[a-z]+\:[0-9]{5})`, address)
	return match
}
