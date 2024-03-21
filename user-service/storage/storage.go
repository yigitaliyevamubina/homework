package storage

import (
	"exam/user-service/pkg/logger"
	mongodb "exam/user-service/storage/mongo"

	// "exam/user-service/storage/postgres"
	"exam/user-service/storage/repo"

	"go.mongodb.org/mongo-driver/mongo"
)

// Storage
type StorageI interface {
	UserService() repo.UserServiceI
}

type storagePg struct {
	userService repo.UserServiceI
}

func New(collection *mongo.Collection, log logger.Logger) StorageI {
	// return &storagePg{userService: postgres.NewUserRepo(db, log)}
	return &storagePg{userService: mongodb.NewUserRepo(collection, log)}
}

func (s *storagePg) UserService() repo.UserServiceI {
	return s.userService
}
