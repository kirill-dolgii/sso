package auth

import (
	"context"

	ssov1 "github.com/kirill-dolgii/protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appId int,
	) (token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)

	IsAdmin(ctx context.Context,
		userID int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const (
	emptyValue = 0
)

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if err := s.validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := s.validateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := s.validateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func (s *serverAPI) validateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "empty email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "empty password")
	}

	if req.AppId == emptyValue {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}
	return nil
}

func (s *serverAPI) validateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "empty email")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "empty password")
	}

	return nil
}

func (s *serverAPI) validateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() < 0 {
		return status.Error(codes.InvalidArgument, "empty user id")
	}

	return nil
}
