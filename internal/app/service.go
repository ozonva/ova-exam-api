package ova_exam_api

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	desc "ova-exam-api/pkg/github.com/ozonva/ova-exam-api/pkg/ova-exam-api"
)

type OvaExamAPI struct {
	desc.UnimplementedUsersServer
}

func NewOvaExamAPI() desc.UsersServer {
	return &OvaExamAPI{
	}
}
func (a *OvaExamAPI)CreateUserV1(ctx context.Context, req *desc.CreateUserV1Request) (*emptypb.Empty, error) {
	log.Print(req)
	return &emptypb.Empty{}, nil
}

func (a *OvaExamAPI) DescribeUserV1(ctx context.Context, req *desc.DescribeUserV1Request) (*desc.UserV1, error) {
	log.Print(req)
	return nil, status.Errorf(codes.Unimplemented, "method DescribeUserV1 not implemented =)")
}
func (a *OvaExamAPI) ListUsersV1(ctx context.Context, req *empty.Empty) (*desc.ListUsersV1Response, error) {
	log.Print(req)
	return nil, status.Errorf(codes.Unimplemented, "method ListUsersV1 not implemented =)")
}
func (a *OvaExamAPI) RemoveUserV1(ctx context.Context, req *desc.RemoveUserV1Request) (*empty.Empty, error) {
	log.Print(req)
	return nil, status.Errorf(codes.Unimplemented, "method RemoveUserV1 not implemented =)")
}