package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	DBHost     string `json:"DBHost"`
	DBPort     string `json:"DBPort"`
	DBUsername string `json:"DBUsername"`
	DBPassword string `json:"DBPassword"`
	DBName     string `json:"DBName"`
}

type Database struct {
	Client *sqlx.DB
}

func (d *Database) GetClient() *sqlx.DB {
	return d.Client
}

func NewDatabase(configPath string) (*Database, error) {
	log.Info("Setting up new database connection")

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("could not open config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("could not decode config file: %w", err)
	}

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUsername,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
		config.DBName,
	)

	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("could not connect to database: %w", err)
	}

	return &Database{
		Client: db,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}
