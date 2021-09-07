package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"ova-exam-api/internal/repo"

	_ "github.com/jackc/pgx/stdlib"
	"google.golang.org/grpc"

	api "ova-exam-api/internal/app"

	desc "ova-exam-api/pkg/github.com/ozonva/ova-exam-api/pkg/ova-exam-api"
)

const (
	grpcPort = "localhost:82"
)

func run() error {
	connectionString := "host=localhost port=5432 user=ova_exam_api_user password=ova_exam_api_password dbname=ova_exam_api sslmode=disable"

	usersRepo := repo.NewRepo(connectionString)

	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to listen: %v", err))
	}

	s := grpc.NewServer()
	desc.RegisterUsersServer(s, api.NewOvaExamAPI(usersRepo))

	if err := s.Serve(listen); err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to serve: %v", err))
	}

	return nil
}


func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	if err := run(); err != nil {
		log.Fatal().Err(err)
	}
}