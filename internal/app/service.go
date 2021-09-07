package ovaexamapi

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"ova-exam-api/internal/domain/entity/user"
	"ova-exam-api/internal/repo"
	desc "ova-exam-api/pkg/github.com/ozonva/ova-exam-api/pkg/ova-exam-api"
)

type OvaExamAPI struct {
	desc.UnimplementedUsersServer
	repo repo.Repo
}

func NewOvaExamAPI(repo repo.Repo) desc.UsersServer {
	return &OvaExamAPI{
		repo: repo,
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

	return &emptypb.Empty{}, nil
}

func (a *OvaExamAPI) DescribeUserV1(ctx context.Context, req *desc.DescribeUserV1Request) (*desc.UserV1, error) {
	log.Log().Int64("UserId", req.UserId).Msg("Запрос на получение из БД пользователя с Id")

	existUser, err := a.repo.DescribeEntity(uint64(req.UserId))
	if err != nil{
		return nil, err
	}
		result := desc.UserV1{
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

	usersResult := make([]*desc.UserV1, 0, len(users))

	for _, us := range users {
		userDto := desc.UserV1 {
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
	
	var result empty.Empty
	return &result, err
}