package main

import (
	"collector/logs"
	"collector/metrics"
	pb "collector/proto"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8010", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewServerClient(conn)
	logs.Watch()
	metrics.CollectAndSendData(client)
}
