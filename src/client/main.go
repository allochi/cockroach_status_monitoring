package main

import (
	"context"
	"log"

	models "../models"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:8899", grpc.WithInsecure())
	if err != nil {
		log.Fatalln("Couldn't dial on localhost:8899")
	}
	defer conn.Close()

	client := models.NewClusterClient(conn)
	status, err := client.GetStatus(context.Background(), &models.Void{})
	if err != nil {
		log.Fatalln("Couldn't call GetStatus", err)
	}
	log.Printf("response: %v", status)
}
