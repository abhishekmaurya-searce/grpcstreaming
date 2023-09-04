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
	unaryResponse, err := client.UnaryExample(context.Background(), &pb.RequestMessage{Message: "Hello Unary"})
	if err != nil {
		log.Fatalf("Unary call failed: %v", err)
	}
	fmt.Printf("%s\n", unaryResponse.Message)
}
