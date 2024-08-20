package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	student "Students-Final-Assignment/Internal/Student"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"github.com/go-playground/validator/v10"
)

type StudentService interface {
	GetStudent(ctx context.Context, ID int64) (student.Student, error)
	PostStudent(ctx context.Context, s student.Student) (student.Student, error)
	UpdateStudent(ctx context.Context, ID int64, s student.Student) (student.Student, error)
	DeleteStudent(ctx context.Context, ID int64) error
	ReadyCheck(ctx context.Context) error
}

func (h *Handler) GetStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := h.Service.GetStudent(r.Context(), id)
	if err != nil {
		if errors.Is(err, student.ErrFetchingStudent) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}
}

type PostStudentRequest struct {
	FirstName   string `json:"fname" validate:"required"`
	LastName    string `json:"lname" validate:"required"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Address     string `json:"address" validate:"required"`
	Gender      string `json:"gender" validate:"required"`
}

func studentFromPostStudentRequest(u PostStudentRequest) student.Student {
	dateOfBirth, _ := time.Parse("2006-01-02", u.DateOfBirth)
	return student.Student{
		Fname:       u.FirstName,
		Lname:       u.LastName,
		DateOfBirth: dateOfBirth,
		Email:       u.Email,
		Address:     u.Address,
		Gender:      u.Gender,
	}
}

func (h *Handler) PostStudent(w http.ResponseWriter, r *http.Request) {
	var postStudentReq PostStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&postStudentReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(postStudentReq); err != nil {
		log.Info(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := studentFromPostStudentRequest(postStudentReq)
	s, err := h.Service.PostStudent(r.Context(), s)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}
}

type UpdateStudentRequest struct {
	FirstName   string    `json:"fname" validate:"required"`
	LastName    string    `json:"lname" validate:"required"`
	DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
	Email       string    `json:"email" validate:"required,email"`
	Address     string    `json:"address" validate:"required"`
	Gender      string    `json:"gender" validate:"required"`
}

func studentFromUpdateStudentRequest(u UpdateStudentRequest) student.Student {
	return student.Student{
		Fname:       u.FirstName,
		Lname:       u.LastName,
		DateOfBirth: u.DateOfBirth,
		Email:       u.Email,
		Address:     u.Address,
		Gender:      u.Gender,
	}
}

func (h *Handler) UpdateStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var updateStudentRequest UpdateStudentRequest
	if err := json.NewDecoder(r.Body).Decode(&updateStudentRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(updateStudentRequest); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := studentFromUpdateStudentRequest(updateStudentRequest)
	s, err = h.Service.UpdateStudent(r.Context(), id, s)
	if err != nil {
		log.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(s); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteStudent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.Service.DeleteStudent(r.Context(), id)
	if err != nil {
		if errors.Is(err, student.ErrDeletingStudent) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "Successfully Deleted"}); err != nil {
		panic(err)
	}
}
