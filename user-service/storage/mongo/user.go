package mongo

import (
	"context"
	pb "exam/user-service/genproto/user-service"
	"exam/user-service/pkg/logger"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type userRepo struct {
	collection *mongo.Collection
	log        logger.Logger
}

func NewUserRepo(collection *mongo.Collection, log logger.Logger) *userRepo {
	return &userRepo{collection: collection, log: log}
}

func (u *userRepo) CreateUser(ctx context.Context, user *pb.User) (*pb.User, error) {
	user.Id = uuid.NewString()
	result, err := u.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	var response pb.User
	filter := bson.M{"_id": result.InsertedID}
	err = u.collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (u *userRepo) GetUserById(ctx context.Context, userId *pb.GetUserId) (*pb.User, error) {
	var response pb.User
	filter := bson.M{"id": userId.UserId}
	err := u.collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	var response pb.User

	filter := bson.M{"_id": req.Id}

	updateReq := bson.M{
		"$set": bson.M{
			"first_name": req.FirstName,
			"last_name":  req.LastName,
			"age":        req.Age,
			"updated_at": time.Now(),
		},
	}

	err := u.collection.FindOneAndUpdate(ctx, filter, updateReq).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (u *userRepo) DeleteUser(ctx context.Context, req *pb.GetUserId) (*pb.Status, error) {
	filter := bson.M{"id": req.UserId}
	_, err := u.collection.DeleteOne(ctx, filter)
	if err != nil {
		return &pb.Status{Success: false}, err
	}

	return &pb.Status{Success: true}, nil
}

func (u *userRepo) ListUsers(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	var response pb.GetListResponse

	reqOptions := options.Find()

	reqOptions.SetSkip(int64((req.Page - 1) * req.Limit))
	reqOptions.SetLimit(int64(req.Limit))

	cursor, err := u.collection.Find(ctx, bson.M{}, reqOptions)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var user pb.User
		err = cursor.Decode(&user)
		if err != nil {
			return nil, err
		}

		response.Count++
		response.Users = append(response.Users, &user)
	}

	return &response, nil
}

func (u *userRepo) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	filter := bson.M{req.Field: req.Data}
	err := u.collection.FindOne(ctx, filter)
	if err != nil {
		return &pb.CheckFieldResponse{Status: false}, err.Err()
	}

	return &pb.CheckFieldResponse{Status: true}, nil
}

func (u *userRepo) Check(ctx context.Context, req *pb.IfExists) (*pb.User, error) {
	var response pb.User
	filter := bson.M{"email": req.Email}
	err := u.collection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (u *userRepo) UpdateRefreshToken(ctx context.Context, req *pb.UpdateRefreshTokenReq) (*pb.Status, error) {
	filter := bson.M{"id": req.UserId}
	updateReq := bson.M{
		"$set": bson.M{
			"refresh_token": req.RefreshToken,
		},
	}

	updateResult, err := u.collection.UpdateOne(ctx, filter, updateReq)
	if err != nil {
		return nil, err
	}
	if updateResult.ModifiedCount == 0 {
		return &pb.Status{Success: false}, nil
	}

	return &pb.Status{Success: true}, nil
}
