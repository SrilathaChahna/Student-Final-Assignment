package database

import (
	"Students-Final-Assignment/Internal/User"
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type SQLUserStore struct {
	Client *sqlx.DB
}

func NewUserStore(db *sqlx.DB) User.UserStore {
	return &SQLUserStore{Client: db}
}

func (s *SQLUserStore) GetUserByUsername(ctx context.Context, username string) (User.User, error) {
	var user User.User
	err := s.Client.GetContext(
		ctx,
		&user,
		`SELECT uid, username, password, email, COALESCE(jwt_token, '') as jwt_token, created_on, updated_on
		FROM users 
		WHERE username = ?`,
		username,
	)
	if err != nil {
		return User.User{}, fmt.Errorf("an error occurred fetching user by username: %w", err)
	}
	return user, nil
}

func (s *SQLUserStore) GetUserByID(ctx context.Context, id int64) (User.User, error) {
	var user User.User
	err := s.Client.GetContext(ctx, &user, "SELECT uid, username, password, email, jwt_token FROM users WHERE uid = ?", id)
	if err != nil {
		return User.User{}, err
	}
	return user, nil
}

func (s *SQLUserStore) CreateUser(ctx context.Context, user User.User) error {
	_, err := s.Client.ExecContext(ctx, "INSERT INTO users (username, password, email) VALUES (?, ?, ?)", user.Username, user.Password, user.Email)
	return err
}

func (s *SQLUserStore) UpdateUser(ctx context.Context, user User.User) error {
	_, err := s.Client.ExecContext(ctx, "UPDATE users SET password = ?, email = ?, jwt_token = ? WHERE uid = ?", user.Password, user.Email, user.JWTToken, user.UID)
	return err
}

func (s *SQLUserStore) DeleteUser(ctx context.Context, id int64) error {
	_, err := s.Client.ExecContext(ctx, "DELETE FROM users WHERE uid = ?", id)
	return err
}

func (s *SQLUserStore) Ping(ctx context.Context) error {
	return s.Client.PingContext(ctx)
}
