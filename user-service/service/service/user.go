package service

import (
	"context"
	pb "exam/user-service/genproto/user-service"
	"exam/user-service/pkg/logger"
	grpcClient "exam/user-service/service/grpc_client"
	"exam/user-service/storage"
)

type UserService struct {
	storage storage.StorageI
	log     logger.Logger
	service grpcClient.IServiceManager
}

// Constructor
func NewUserService(storage storage.StorageI, log logger.Logger, service grpcClient.IServiceManager) *UserService {
	return &UserService{
		storage: storage,
		log:     log,
		service: service,
	}
}

func (c *UserService) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	return c.storage.UserService().CreateUser(ctx, req)
}

func (c *UserService) GetUserById(ctx context.Context, req *pb.GetUserId) (*pb.User, error) {
	return c.storage.UserService().GetUserById(ctx, req)
}

func (c *UserService) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	return c.storage.UserService().UpdateUser(ctx, req)
}

func (c *UserService) DeleteUser(ctx context.Context, req *pb.GetUserId) (*pb.Status, error) {
	return c.storage.UserService().DeleteUser(ctx, req)
}

func (c *UserService) ListUsers(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	return c.storage.UserService().ListUsers(ctx, req)
}

func (c *UserService) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	return c.storage.UserService().CheckField(ctx, req)
}

func (c *UserService) Check(ctx context.Context, req *pb.IfExists) (*pb.User, error) {
	return c.storage.UserService().Check(ctx, req)
}

func (c *UserService) UpdateRefreshToken(ctx context.Context, req *pb.UpdateRefreshTokenReq) (*pb.Status, error) {
	return c.storage.UserService().UpdateRefreshToken(ctx, req)
}
