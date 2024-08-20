package http

import (
	"Students-Final-Assignment/Internal/User"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Handler struct {
	Router      *mux.Router
	Service     StudentService
	Server      *http.Server
	UserService *User.Service
}

type Response struct {
	Message string `json:"message"`
}

func NewHandler(service StudentService, userService *User.Service) *Handler {
	log.Info("setting up our handler")
	h := &Handler{
		Service:     service,
		UserService: userService,
	}

	h.Router = mux.NewRouter()
	h.Router.Use(JSONMiddleware)
	h.Router.Use(LoggingMiddleware)
	h.Router.Use(TimeoutMiddleware)
	h.mapRoutes()

	h.Server = &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      h.Router,
	}
	return h
}

func (h *Handler) mapRoutes() {
	h.Router.HandleFunc("/alive", h.AliveCheck).Methods("GET")
	h.Router.HandleFunc("/ready", h.ReadyCheck).Methods("GET")
	h.Router.HandleFunc("/api/v1/student", JWTAuth(h.PostStudent)).Methods("POST")
	h.Router.HandleFunc("/api/v1/student/{id}", JWTAuth(h.GetStudent)).Methods("GET")
	h.Router.HandleFunc("/api/v1/student/{id}", JWTAuth(h.UpdateStudent)).Methods("PUT")
	h.Router.HandleFunc("/api/v1/student/{id}", JWTAuth(h.DeleteStudent)).Methods("DELETE")
	h.Router.HandleFunc("/api/v1/login", h.Login).Methods("POST")
	h.Router.HandleFunc("/api/v1/register", h.Register).Methods("POST")
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var creds User.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	token, err := h.UserService.Login(creds.Username, creds.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var creds User.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.UserService.Register(creds.Username, creds.Password, creds.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) AliveCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Alive!"}); err != nil {
		panic(err)
	}
}

func (h *Handler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	if err := h.Service.ReadyCheck(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(Response{Message: "I am Ready!"}); err != nil {
		panic(err)
	}
}

func (h *Handler) Serve() error {
	go func() {
		if err := h.Server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	h.Server.Shutdown(ctx)

	log.Println("shutting down gracefully")
	return nil
}
