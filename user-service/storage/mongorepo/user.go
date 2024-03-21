package mongorepo

import (
	"context"
	pb "exam/user-service/genproto/user-service"
)

// UserService interface
type UserServiceI interface {
	CreateUser(ctx context.Context, req *pb.User) (*pb.User, error)
	GetUserById(ctx context.Context, req *pb.GetUserId) (*pb.User, error)
	UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error)
	DeleteUser(ctx context.Context, req *pb.GetUserId) (*pb.Status, error)
	ListUsers(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error)
	CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error)
	Check(ctx context.Context, req *pb.IfExists) (*pb.User, error)
	UpdateRefreshToken(ctx context.Context, req *pb.UpdateRefreshTokenReq) (*pb.Status, error)
}
