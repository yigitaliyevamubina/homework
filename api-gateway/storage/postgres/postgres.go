package postgres

import (
	"errors"
	"exam/api-gateway/api/handlers/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type adminRepo struct {
	db *sqlx.DB
}

func NewAdminRepo(db *sqlx.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (r *adminRepo) Create(admin *models.AdminResp) error {
	query := `INSERT INTO admins(id, full_name, age, username, email, password, role, refresh_token) 
								VALUES($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(query, admin.Id,
		admin.FullName,
		admin.Age,
		admin.UserName,
		admin.Email,
		admin.Password,
		admin.Role,
		admin.RefreshToken)

	return err
}

func (r *adminRepo) Delete(userName, password string) error {
	query := `DELETE FROM admins WHERE username = $1`
	result, err := r.db.Exec(query, userName)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		fmt.Println("error")
		return errors.New("no rows were deleted")
	}

	return nil
}

func (r *adminRepo) Get(userName string) (string, string, bool, error) {
	query := `SELECT COUNT(1), password, role
	FROM admins GROUP by username, password, role
	HAVING username = $1;
	`
	var status int
	var password string
	var role string
	result := r.db.QueryRow(query, userName)
	if err := result.Scan(&status, &password, &role); err != nil {
		return "", "", false, nil
	}

	return role, password, status == 1, nil
}
