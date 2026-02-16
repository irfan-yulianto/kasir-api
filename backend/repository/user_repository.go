package repository

import (
	"database/sql"

	"kasir-api/model"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.QueryRow(
		"INSERT INTO users (username, password, role) VALUES ($1, $2, $3) RETURNING id, created_at",
		user.Username, user.Password, user.Role,
	).Scan(&user.ID, &user.CreatedAt)
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var u model.User
	err := r.db.QueryRow(
		"SELECT id, username, password, role, created_at FROM users WHERE username = $1",
		username,
	).Scan(&u.ID, &u.Username, &u.Password, &u.Role, &u.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}
