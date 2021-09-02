package main

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"ova-exam-api/internal/repo"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"

	api "ova-exam-api/internal/app"

	desc "ova-exam-api/pkg/github.com/ozonva/ova-exam-api/pkg/ova-exam-api"
)

const (
	grpcPort = ":82"
	grpcServerEndpoint = "localhost:82"
)

func run() error {
	connectString := "host=localhost port=5432 user=ova_exam_api_user password=ova_exam_api_password dbname=ova_exam_api sslmode=disable"

	db, err := sqlx.Connect("pgx", connectString)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to load driver: %v", err))
	}
	defer db.Close()
	ctx := context.Background()

	dbPingErr := db.PingContext(ctx)
	if dbPingErr != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to connect to db: %v", err))
	}

	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to listen: %v", err))
	}

	usersRepo := repo.NewRepo(db)

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