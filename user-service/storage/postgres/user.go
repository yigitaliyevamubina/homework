package postgres

import (
	"context"
	pb "exam/user-service/genproto/user-service"
	"exam/user-service/pkg/db"
	"exam/user-service/pkg/logger"
	"exam/user-service/storage/repo"
	"github.com/Masterminds/squirrel"
	"time"
)

type userRepo struct {
	db  *db.Postgres
	log logger.Logger
}

// Constructor
func NewUserRepo(db *db.Postgres, log logger.Logger) repo.UserServiceI {
	return &userRepo{
		db:  db,
		log: log,
	}
}

func (u *userRepo) CreateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	query := u.db.Builder.Insert("users").
		Columns(`
		id, first_name, last_name, age, email, password, refresh_token
		`).
		Values(
			req.Id, req.FirstName, req.LastName,
			req.Age, req.Email, req.Password,
			req.RefreshToken,
		).
		Suffix("RETURNING created_at")

	err := query.RunWith(u.db.DB).QueryRow().Scan(&req.CreatedAt)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *userRepo) GetUserById(ctx context.Context, req *pb.GetUserId) (*pb.User, error) {
	respUser := &pb.User{}

	query := u.db.Builder.Select(`
		id, first_name, last_name, age, email, password, refresh_token, created_at
	`).From("users").Where(squirrel.Eq{"id": req.UserId})

	err := query.RunWith(u.db.DB).QueryRow().Scan(
		&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Age,
		&respUser.Email,
		&respUser.Password,
		&respUser.RefreshToken,
		&respUser.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return respUser, nil
}

func (u *userRepo) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	var (
		updateMap = make(map[string]interface{})
		where     = squirrel.And{squirrel.Eq{"id": req.Id}}
	)

	updateMap["first_name"] = req.FirstName
	updateMap["last_name"] = req.LastName
	updateMap["age"] = req.Age
	updateMap["updated_at"] = time.Now()

	query := u.db.Builder.Update("users").SetMap(updateMap).
		Where(where).
		Suffix("RETURNING updated_at, created_at")

	err := query.RunWith(u.db.DB).QueryRow().Scan(
		&req.UpdatedAt, &req.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (u *userRepo) DeleteUser(ctx context.Context, req *pb.GetUserId) (*pb.Status, error) {
	query := u.db.Builder.Delete("users").Where(squirrel.Eq{"id": req.UserId})
	_, err := query.RunWith(u.db.DB).Exec()
	if err != nil {
		return &pb.Status{
			Success: false,
		}, err
	}

	return &pb.Status{
		Success: true,
	}, nil
}

func (u *userRepo) ListUsers(ctx context.Context, req *pb.GetListRequest) (*pb.GetListResponse, error) {
	var (
		respUsers = &pb.GetListResponse{Count: 0}
	)

	query := u.db.Builder.Select(
		`id, first_name, last_name, age, email, password, refresh_token
	`).From("users")

	query = query.Offset(uint64((req.Page - 1) * req.Limit)).Limit(uint64(req.Limit))

	rows, err := query.RunWith(u.db.DB).Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		respUser := &pb.User{}
		err = rows.Scan(
			&respUser.Id,
			&respUser.FirstName,
			&respUser.LastName,
			&respUser.Age,
			&respUser.Email,
			&respUser.Password,
			&respUser.RefreshToken,
		)
		if err != nil {
			return nil, err
		}
		respUsers.Users = append(respUsers.Users, respUser)
		respUsers.Count++
	}

	return respUsers, nil
}

func (u *userRepo) CheckField(ctx context.Context, req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	var (
		response = &pb.CheckFieldResponse{}
	)
	var resp int
	num := u.db.Builder.Select("count(1)").From("users").Where(squirrel.Eq{req.Field: req.Data})

	err := num.RunWith(u.db.DB).Scan(&resp)
	if err != nil {
		response.Status = false
		return response, err
	}
	if resp == 1 {
		response.Status = true
	} else if resp == 0 {
		response.Status = false
	}

	return response, nil
}

func (u *userRepo) Check(ctx context.Context, req *pb.IfExists) (*pb.User, error) {
	respUser := &pb.User{}

	query := u.db.Builder.Select(`
		id, first_name, last_name, age, email, password, refresh_token, created_at
	`).From("users").Where(squirrel.Eq{"email": req.Email})

	err := query.RunWith(u.db.DB).QueryRow().Scan(
		&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
		&respUser.Age,
		&respUser.Email,
		&respUser.Password,
		&respUser.RefreshToken,
		&respUser.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return respUser, nil
}

func (u *userRepo) UpdateRefreshToken(ctx context.Context, req *pb.UpdateRefreshTokenReq) (*pb.Status, error) {
	var (
		updateMap = make(map[string]interface{})
		where     = squirrel.And{squirrel.Eq{"id": req.UserId}}
	)

	updateMap["refresh_token"] = req.RefreshToken

	query := u.db.Builder.Update("users").SetMap(updateMap).
		Where(where).
		Suffix("RETURNING updated_at, created_at")

	_, err := query.RunWith(u.db.DB).Exec()
	if err != nil {
		return &pb.Status{
			Success: false,
		}, err
	}

	return &pb.Status{
		Success: true,
	}, nil
}
