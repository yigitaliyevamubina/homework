package postgres

// import (
// 	"context"
// 	"exam/user-service/config"
// 	pb "exam/user-service/genproto/user-service"
// 	db2 "exam/user-service/pkg/db"
// 	"exam/user-service/pkg/logger"
// 	"exam/user-service/storage/repo"
// 	"github.com/brianvoe/gofakeit/v6"
// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/suite"
// 	"testing"
// 	"time"
// )

// type UserTestSuite struct {
// 	suite.Suite
// 	CleanupFunc func()
// 	Repository  repo.UserServiceI
// }

// func (u *UserTestSuite) SetupSuite() {
// 	db, _ := db2.New(*config.Load())
// 	u.Repository = NewUserRepo(db, logger.New("", ""))
// 	u.CleanupFunc = db.Close
// }

// func (u *UserTestSuite) TestPositionCRUD() {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(7))
// 	defer cancel()
// 	//Create user
// 	id := uuid.New().String()
// 	user := &pb.User{
// 		Id:        id,
// 		FirstName: gofakeit.FirstName(),
// 		LastName:  gofakeit.LastName(),
// 		Age:       10,
// 		Email:     gofakeit.Email(),
// 		Password:  gofakeit.Phrase(),
// 	}

// 	createResp, err := u.Repository.CreateUser(ctx, user)
// 	u.Suite.NoError(err)
// 	u.Suite.NotNil(createResp)

// 	//Get user
// 	userId := &pb.GetUserId{
// 		UserId: id,
// 	}
// 	getResp, err := u.Repository.GetUserById(ctx, userId)
// 	u.Suite.NoError(err)
// 	u.Suite.NotNil(getResp)
// 	u.Suite.Equal(getResp.Age, user.Age)
// 	u.Suite.Equal(getResp.LastName, user.LastName)
// 	u.Suite.Equal(getResp.FirstName, user.FirstName)
// 	u.Suite.Equal(getResp.Email, user.Email)

// 	//List users
// 	listResp, err := u.Repository.ListUsers(ctx, &pb.GetListRequest{
// 		Page:  1,
// 		Limit: 10,
// 	})
// 	u.Suite.NoError(err)
// 	u.Suite.NotNil(listResp)

// 	//Update user
// 	updatedName := gofakeit.FirstName()
// 	user.FirstName = updatedName
// 	updatedAge := int64(13)
// 	user.Age = updatedAge
// 	user.Id = userId.UserId
// 	updateResp, err := u.Repository.UpdateUser(ctx, user)
// 	u.Suite.NoError(err)
// 	u.Suite.NotNil(updateResp)
// 	u.Suite.Equal(updatedName, updateResp.FirstName)
// 	u.Suite.Equal(updatedAge, updateResp.Age)

// 	//CheckField
// 	checkResp, err := u.Repository.CheckField(ctx, &pb.CheckFieldRequest{
// 		Field: "email",
// 		Data:  user.Email,
// 	})
// 	u.Suite.NoError(err)
// 	u.Suite.NotNil(checkResp)
// 	u.Suite.Equal(checkResp.Status, true)

// 	//Delete user
// 	_, err = u.Repository.DeleteUser(ctx, userId)
// 	u.Suite.NoError(err)
// }

// func (u *UserTestSuite) TearDownSuite() {
// 	u.CleanupFunc()
// }

// func TestUserRepository(t *testing.T) {
// 	suite.Run(t, new(UserTestSuite))
// }
