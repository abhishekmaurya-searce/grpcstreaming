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
	bidiStream, err := client.BidirectionalStreamingExample(context.Background())
	if err != nil {
		log.Fatalf("Bidirectional Streaming call failed: %v", err)
	}
	go func() {
		for i := 0; i < 5; i++ {
			req := &pb.RequestMessage{Message: fmt.Sprintf("Bidirectional Streaming %d", i+1)}
			if err := bidiStream.Send(req); err != nil {
				log.Fatalf("Bidirectional Streaming Send failed: %v", err)
			}
			time.Sleep(1 * time.Second)
		}
		bidiStream.CloseSend()
	}()
	for {
		resp, err := bidiStream.Recv()
		if err != nil {
			break
		}
		fmt.Printf("Bidirectional Streaming Response: %s\n", resp.Message)
	}
}
