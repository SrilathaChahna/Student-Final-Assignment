package database

import (
	"context"
	"fmt"
	"time"

	"Students-Final-Assignment/Internal/Student"

	"github.com/jmoiron/sqlx"
)

type StudentRow struct {
	ID          int64     `db:"id"`
	FName       string    `db:"fname"`
	LName       string    `db:"lname"`
	DateOfBirth time.Time `db:"date_of_birth"`
	Email       string    `db:"email"`
	Address     string    `db:"address"`
	Gender      string    `db:"gender"`
	CreatedBy   string    `db:"created_by"`
	CreatedOn   time.Time `db:"created_on"`
}

type SQLStudentStore struct {
	Client *sqlx.DB
}

func NewStudentStore(db *sqlx.DB) Student.StudentStore {
	return &SQLStudentStore{Client: db}
}

func (s *SQLStudentStore) Ping(ctx context.Context) error {
	return s.Client.PingContext(ctx)
}

func convertStudentRowToStudent(row StudentRow) Student.Student {
	return Student.Student{
		ID:          row.ID,
		Fname:       row.FName,
		Lname:       row.LName,
		DateOfBirth: row.DateOfBirth,
		Email:       row.Email,
		Address:     row.Address,
		Gender:      row.Gender,
		CreatedBy:   row.CreatedBy,
		CreatedOn:   row.CreatedOn,
	}
}

func (s *SQLStudentStore) GetStudent(ctx context.Context, id int64) (Student.Student, error) {
	var row StudentRow
	err := s.Client.GetContext(
		ctx,
		&row,
		`SELECT id, fname, lname, date_of_birth, email, address, gender, created_by, created_on
		FROM students 
		WHERE id = ?`,
		id,
	)
	if err != nil {
		return Student.Student{}, fmt.Errorf("an error occurred fetching a student by id: %w", err)
	}
	return convertStudentRowToStudent(row), nil
}

func (s *SQLStudentStore) PostStudent(ctx context.Context, st Student.Student) (Student.Student, error) {
	_, err := s.Client.ExecContext(
		ctx,
		`INSERT INTO students (fname, lname, date_of_birth, email, address, gender, created_by, created_on) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		st.Fname, st.Lname, st.DateOfBirth, st.Email, st.Address, st.Gender, "admin", time.Now(),
	)
	if err != nil {
		return Student.Student{}, fmt.Errorf("failed to insert student: %w", err)
	}
	return st, nil
}

func (s *SQLStudentStore) UpdateStudent(ctx context.Context, id int64, st Student.Student) (Student.Student, error) {
	_, err := s.Client.ExecContext(
		ctx,
		`UPDATE students SET fname = ?, lname = ?, date_of_birth = ?, email = ?, address = ?, gender = ?, updated_by = ?, updated_on = ? WHERE id = ?`,
		st.Fname, st.Lname, st.DateOfBirth, st.Email, st.Address, st.Gender, "admin", time.Now(), id,
	)
	if err != nil {
		return Student.Student{}, fmt.Errorf("failed to update student: %w", err)
	}
	st.ID = id
	return st, nil
}

func (s *SQLStudentStore) DeleteStudent(ctx context.Context, id int64) error {
	_, err := s.Client.ExecContext(
		ctx,
		`DELETE FROM students WHERE id = ?`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete student from the database: %w", err)
	}
	return nil
}
