package main

import (
	"context"
	"fmt"
	"log"
	pb "streaming/proto"
	"time"

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
	clientStream, err := client.ClientStreamingExample(context.Background())
	if err != nil {
		log.Fatalf("Client Streaming call failed: %v", err)
	}
	for i := 0; i < 5; i++ {
		req := &pb.RequestMessage{Message: fmt.Sprintf("Client Streaming %d", i+1)}
		if err := clientStream.Send(req); err != nil {
			log.Fatalf("Client Streaming Send failed: %v", err)
		}
		time.Sleep(1 * time.Second)
	}
	message, err := clientStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error in response %v \n", err)
	}
	fmt.Printf("Client Streaming Response: %s\n", message)

}
