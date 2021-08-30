package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	api "ova-exam-api/internal/app"

	desc "ova-exam-api/pkg/github.com/ozonva/ova-exam-api/pkg/ova-exam-api"
)

const (
	grpcPort = ":82"
	grpcServerEndpoint = "localhost:82"
)

func run() error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	desc.RegisterUsersServer(s, api.NewOvaExamAPI())

	if err := s.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}


func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}