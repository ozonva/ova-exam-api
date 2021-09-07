package ovaexamapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	opentracing "github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/rs/zerolog/log"
	kafka "github.com/segmentio/kafka-go"
	"google.golang.org/protobuf/types/known/emptypb"
	"ova-exam-api/internal/domain/entity/user"
	"ova-exam-api/internal/repo"
	"ova-exam-api/internal/utils"
	desc "ova-exam-api/pkg/github.com/ozonva/ova-exam-api/pkg/ova-exam-api"
	"time"
)

type OvaExamAPI struct {
	desc.UnimplementedUsersServer
	repo repo.Repo
	conn *kafka.Conn
	usersSaved prometheus.Counter
}


func NewOvaExamAPI(repo repo.Repo, usersSaved *prometheus.Counter) desc.UsersServer {
	topic := "test-topic"
	partition := 0

	newConn, err := kafka.DialLeader(context.Background(), "tcp", "127.0.0.1:9092", topic, partition)
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf("failed to dial leader: %v", err))
	}

	return &OvaExamAPI{
		repo:                     repo,
		conn:                     newConn,
		usersSaved:               *usersSaved,
	}
}
func (a *OvaExamAPI)CreateUserV1(ctx context.Context, req *desc.CreateUserV1Request) (*emptypb.Empty, error) {
	log.Log().Str("Email", req.Email).Msg("Запрос на создание пользователя")

	newUser := user.User{
		Email:    req.Email,
		Password: req.Password,
	}
	users := []user.User{newUser}
	err := a.repo.AddEntities(users)
	if err != nil{
		return nil, err
	}

	serializedMessage, marshallErr := json.Marshal(users)
	if marshallErr != nil {
		return &emptypb.Empty{}, marshallErr
	}
	kafkaErr := a.sentToKafka(serializedMessage)
	if kafkaErr != nil {
		return &emptypb.Empty{}, kafkaErr
	}

	a.usersSaved.Inc()
	return &emptypb.Empty{}, nil
}

func (a *OvaExamAPI) MultiCreateUserV1(ctx context.Context, req *desc.MultyCreateUserV1Request) (*emptypb.Empty, error) {
	log.Log().Int("Count", len(req.Users)).Msg("Запрос на создание пользователей")

	span, ctx := opentracing.StartSpanFromContext(ctx, "operation_name")
	defer span.Finish()

	for _, batch := range utils.SplitUsersV1RequestToBulks(req.Users, 2) {
		usersResult := make([]user.User, 0, len(req.Users))

		for _, us := range batch {
			userDto := user.User{
				Email:    us.Email,
				Password: us.Password,
			}
			usersResult = append(usersResult, userDto)
		}
		err := a.repo.AddEntities(usersResult)
		if err != nil {
			return nil, err
		}
		a.usersSaved.Add(float64(len(usersResult)))
		serializedMessage, marshallErr := json.Marshal(usersResult)
		if marshallErr != nil {
			return &emptypb.Empty{}, marshallErr
		}

		opentracing.StartSpan(
			"operation_name",
			opentracing.ChildOf(span.Context())).
			LogFields(otlog.Int("sendSizeInBytes", len(serializedMessage)))

		kafkaErr := a.sentToKafka(serializedMessage)
		if kafkaErr != nil {
			return &emptypb.Empty{}, kafkaErr
		}
	}

	return &emptypb.Empty{}, nil
}

func (a *OvaExamAPI) DescribeUserV1(ctx context.Context, req *desc.DescribeUserV1Request) (*desc.UserV1Response, error) {
	log.Log().Int64("UserId", req.UserId).Msg("Запрос на получение из БД пользователя с Id")

	existUser, err := a.repo.DescribeEntity(uint64(req.UserId))
	if err != nil{
		return nil, err
	}
		result := desc.UserV1Response{
			UserId:   int64(existUser.UserId),
			Email:    existUser.Email,
			Password: existUser.Password,
		}

	return &result, nil
}

func (a *OvaExamAPI) ListUsersV1(ctx context.Context, req *empty.Empty) (*desc.ListUsersV1Response, error) {
	log.Log().Msg("Запрос на получение всех пользователей из БД")
	maxUInt64 := ^uint64(0)
	limit := maxUInt64 >> 1
	users, err := a.repo.ListEntities(limit, 0)
	if err != nil {
		return nil, err
	}

	usersResult := make([]*desc.UserV1Response, 0, len(users))

	for _, us := range users {
		userDto := desc.UserV1Response {
			UserId:   int64(us.UserId),
			Email:    us.Email,
			Password: us.Password,
		}

		usersResult = append(usersResult, &userDto)
	}

	result := &desc.ListUsersV1Response{
		Users: usersResult,
	}

	return result, nil
}

func (a *OvaExamAPI) RemoveUserV1(ctx context.Context, req *desc.RemoveUserV1Request) (*empty.Empty, error) {
	log.Debug().Int64("UserId", req.UserId).Msg("Запрос на удаление из БД пользователя")
	err := a.repo.RemoveEntity(uint64(req.UserId))
	if err != nil{
		return nil, err
	}

	serializedMessage, marshallErr := json.Marshal(req)
	if marshallErr != nil {
		return &emptypb.Empty{}, marshallErr
	}
	kafkaErr := a.sentToKafka(serializedMessage)
	if kafkaErr != nil {
		return &emptypb.Empty{}, kafkaErr
	}

	return & empty.Empty{}, err
}

func (a *OvaExamAPI)UpdateUserV1(ctx context.Context, req *desc.UpdateUserV1Request) (*emptypb.Empty, error) {
	log.Log().Int64("UserId", req.UserId).Msg("Запрос на изменение пользователя")

	updateUser := user.User{
		UserId:   uint64(req.UserId),
		Email:    req.Email,
		Password: req.Password,
	}

	err := a.repo.UpdateEntity(updateUser)
	if err != nil{
		return nil, err
	}

	serializedMessage, marshallErr := json.Marshal(req)
	if marshallErr != nil {
		return &emptypb.Empty{}, marshallErr
	}
	kafkaErr := a.sentToKafka(serializedMessage)
	if kafkaErr != nil {
		return &emptypb.Empty{}, kafkaErr
	}

	a.usersSaved.Inc()
	return &emptypb.Empty{}, nil
}

func (a *OvaExamAPI) sentToKafka(message []byte) error {
	err1 := a.conn.SetWriteDeadline(time.Now().Add(10*time.Second))
	if err1 != nil {
		return err1
	}
	_, err2 := a.conn.WriteMessages(kafka.Message{Value: message})

	if err2 != nil {
		return err2
	}

	return nil
}