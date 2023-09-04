package main

import (
	"context"
	"fmt"
	"log"
	pb "streaming/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = "0.0.0.0:50051"

func main() {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	client := pb.NewStreamingServiceClient(conn)
	serverStream, err := client.ServerStreamingExample(context.Background(), &pb.RequestMessage{Message: "Hello Server Streaming"})
	if err != nil {
		log.Fatalf("Server Streaming call failed: %v", err)
	}
	for {
		resp, err := serverStream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("Server Streaming Response: %s\n", resp.Message)
	}
}
