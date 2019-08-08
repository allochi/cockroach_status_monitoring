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

func main() {
	var entryAddress string
	flag.StringVar(&entryAddress, "address", "", "entry host http address in the format `address:port` or left empty for `localhost:8080")
	var duration = flag.Int("duration", 5, "duration between each status update")
	flag.Parse()

	if entryAddress != "" {
		if !isValidAddress(entryAddress) {
			log.Fatalf("host address must be in the format `address:port` or left empty for `localhost:8080`\n")
		}
	} else {
		entryAddress = "localhost:8080"
	}

	ch = models.ClusterHealth{
		EntryAddress:   entryAddress,
		StatusProvider: providers.HTTPStatusProvider{Address: entryAddress},
	}

	// Update cluster status periodically
	ticker := time.NewTicker(time.Duration(*duration) * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			ch.Update()
			// fmt.Printf("%#v\n", ch)
			// for _, node := range ch.Nodes {
			// 	content, _ := json.Marshal(node)
			// 	log.Println(string(content))
			// }
		}
	}()

	// gRPC server
	srv := grpc.NewServer()

	var statusService StatusService
	models.RegisterHealthServiceServer(srv, statusService)

	ln, err := net.Listen("tcp", ":8899")
	if err != nil {
		log.Fatalln("Couldn't listen on localhost:8899")
	}
	log.Fatal(srv.Serve(ln))

}

// StatusService gRPC service
type StatusService struct{}

// GetStatus gRPC implementation
func (s StatusService) GetStatus(ctx context.Context, void *models.Void) (*models.HealthResponse, error) {
	response := &models.HealthResponse{
		TotalNodes: int64(len(ch.Nodes)),
	}

	var nodes []*models.Node
	for i, node := range ch.Nodes {
		nodes = append(nodes, &ch.Nodes[i])

		if node.IsLive {
			response.TotalNodesLive++
		}

		if node.IsAvailable {
			response.TotalNodesAvailable++
		}

		if node.IsLowInMemory {
			response.TotalNodesLowMemory++
		}

		response.ClusterUnavailable = float64(response.TotalNodesLive)/float64(response.TotalNodes) < 0.5
		response.UpdatedAt = ch.UpdatedAt.UnixNano()
	}
	response.Nodes = nodes
	return response, nil
}

func isValidAddress(address string) bool {
	match, _ := regexp.MatchString(`([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\:[0-9]{5})|(localhost\:[0-9]{5})|([a-z]+\.[a-z]+\.[a-z]+\:[0-9]{5})`, address)
	return match
}
