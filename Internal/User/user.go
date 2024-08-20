package User

import (
	"context"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UID       int64     `db:"uid" json:"uid"`
	Username  string    `db:"username" json:"username"`
	Password  string    `db:"password" json:"password"`
	Email     string    `db:"email" json:"email"`
	JWTToken  *string   `db:"jwt_token" json:"jwt_token"`
	CreatedOn time.Time `db:"created_on" json:"created_on"`
	UpdatedOn time.Time `db:"updated_on" json:"updated_on"`
}

type UserStore interface {
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserByID(ctx context.Context, id int64) (User, error)
	CreateUser(ctx context.Context, user User) error
	UpdateUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id int64) error
	Ping(ctx context.Context) error
}

type Service struct {
	store UserStore
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func NewService(store UserStore) *Service {
	return &Service{store: store}
}

func (s *Service) Login(username, password string) (string, error) {
	user, err := s.store.GetUserByUsername(context.Background(), username)
	if err != nil {
		fmt.Println("Error fetching user:", err)
		return "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		fmt.Println("Invalid password:", err)
		return "", fmt.Errorf("invalid password")
	}

	mySigningKey := []byte("missionimpossible")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": user.UID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("Error generating token:", err)
		return "", err
	}
	user.JWTToken = &tokenString
	err = s.store.UpdateUser(context.Background(), user)
	if err != nil {
		fmt.Println("Error updating user with token:", err)
		return "", err
	}

	return tokenString, nil
}

func (s *Service) Register(username, password, email string) error {
	_, err := s.store.GetUserByUsername(context.Background(), username)
	if err == nil {
		return fmt.Errorf("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Username:  username,
		Password:  string(hashedPassword),
		Email:     email,
		CreatedOn: time.Now(),
	}
	return s.store.CreateUser(context.Background(), user)
}
