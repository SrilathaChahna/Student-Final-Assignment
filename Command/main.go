package main

import (
	database "Students-Final-Assignment/Internal/Database"
	transportHTTP "Students-Final-Assignment/Internal/Services/http"
	"Students-Final-Assignment/Internal/Student"
	"Students-Final-Assignment/Internal/User"

	"go.uber.org/zap"
)

func Run() error {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	logger.Info("Setting Up Our APP")

	var dbErr error
	db, dbErr := database.NewDatabase("D:/training/GoLang/Students-Final-Assignment/Internal/Database/config.json")
	if dbErr != nil {
		logger.Error("failed to setup connection to the database", zap.Error(dbErr))
		return dbErr
	}

	studentStore := database.NewStudentStore(db.GetClient())
	userStore := database.NewUserStore(db.GetClient())

	studentService := Student.NewService(studentStore)
	userService := User.NewService(userStore)
	handler := transportHTTP.NewHandler(studentService, userService)

	if serveErr := handler.Serve(); serveErr != nil {
		logger.Error("failed to gracefully serve our application", zap.Error(serveErr))
		return serveErr
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		logger, _ := zap.NewProduction()
		defer logger.Sync()
		logger.Error("Error starting up our REST API", zap.Error(err))
	}
}
