package Student

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrFetchingStudent = errors.New("could not fetch Student by ID")
	ErrUpdatingStudent = errors.New("could not update Student")
	ErrNoStudentFound  = errors.New("no Student found")
	ErrDeletingStudent = errors.New("could not delete Student")
	ErrNotImplemented  = errors.New("not implemented")
)

type Student struct {
	ID          int64     `json:"id"`
	Fname       string    `json:"fname"`
	Lname       string    `json:"lname"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Email       string    `json:"email"`
	Address     string    `json:"address"`
	Gender      string    `json:"gender"`
	CreatedBy   string    `json:"created_by"`
	CreatedOn   time.Time `json:"created_on"`
	UpdatedBy   int       `json:"updated_by"`
	UpdatedOn   time.Time `json:"updated_on"`
}

type StudentStore interface {
	GetStudent(context.Context, int64) (Student, error)
	PostStudent(context.Context, Student) (Student, error)
	UpdateStudent(context.Context, int64, Student) (Student, error)
	DeleteStudent(context.Context, int64) error
	Ping(context.Context) error
}

type Service struct {
	Store StudentStore
}

func NewService(store StudentStore) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetStudent(ctx context.Context, ID int64) (Student, error) {
	cmt, err := s.Store.GetStudent(ctx, ID)
	if err != nil {
		log.Errorf("an error occured fetching the Student: %s", err.Error())
		return Student{}, ErrFetchingStudent
	}
	return cmt, nil
}

func (s *Service) PostStudent(ctx context.Context, cmt Student) (Student, error) {
	cmt, err := s.Store.PostStudent(ctx, cmt)
	if err != nil {
		log.Errorf("an error occurred adding the Student: %s", err.Error())
	}
	return cmt, nil
}

func (s *Service) UpdateStudent(
	ctx context.Context, ID int64, newStudent Student,
) (Student, error) {
	cmt, err := s.Store.UpdateStudent(ctx, ID, newStudent)
	if err != nil {
		log.Errorf("an error occurred updating the Student: %s", err.Error())
	}
	return cmt, nil
}

func (s *Service) DeleteStudent(ctx context.Context, ID int64) error {
	return s.Store.DeleteStudent(ctx, ID)
}

func (s *Service) ReadyCheck(ctx context.Context) error {
	log.Info("Checking readiness")
	return s.Store.Ping(ctx)
}
