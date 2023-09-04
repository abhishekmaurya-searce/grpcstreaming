package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	pb "streaming/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	pb.StreamingServiceServer
}

func (s *Server) UnaryExample(ctx context.Context, req *pb.RequestMessage) (*pb.ResponseMessage, error) {
	return &pb.ResponseMessage{Message: "Unary Response: " + req.Message}, nil
}
func (s *Server) ServerStreamingExample(req *pb.RequestMessage, stream pb.StreamingService_ServerStreamingExampleServer) error {
	for i := 0; i < 5; i++ {
		resp := &pb.ResponseMessage{Message: fmt.Sprintf("Server Streaming Response %d: %s", i+1, req.Message)}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
	return nil
}
func (s *Server) ClientStreamingExample(stream pb.StreamingService_ClientStreamingExampleServer) error {
	var messages []string
	for {
		req, err := stream.Recv()
		if err != nil {
			break
		}
		messages = append(messages, req.Message)
	}
	return stream.SendAndClose(&pb.ResponseMessage{Message: "Client Streaming Response: " + fmt.Sprint(messages)})
}
func (s *Server) BidirectionalStreamingExample(stream pb.StreamingService_BidirectionalStreamingExampleServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		resp := &pb.ResponseMessage{Message: "Bidirectional Streaming Response: " + req.Message}
		if err := stream.Send(resp); err != nil {
			return err
		}
	}
}
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	// with Credantials
	serverCert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to load server certificates: %v", err)
	}

	// Load CA certificate
	caCert, err := os.ReadFile("ca.crt")
	if err != nil {
		log.Fatalf("Failed to load CA certificate: %v", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(caCert) {
		log.Fatalf("Failed to append CA certificate")
	}

	// Create TLS credentials for the server
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	})

	// Create a gRPC server with TLS
	s := grpc.NewServer(grpc.Creds(creds))
	// s := grpc.NewServer()
	pb.RegisterStreamingServiceServer(s, &Server{})

	fmt.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
