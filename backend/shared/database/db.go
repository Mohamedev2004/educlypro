package database

import (
	"database/sql"
	"educlypro/config"
	"fmt"
	"log"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib" // Standard Postgres SQL driver for DB creation
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	MainDB *gorm.DB
	LogsDB *gorm.DB
)

func Connect() {
	dbConfig := config.Cfg.DB
	user := dbConfig.User
	pass := dbConfig.Pass
	host := dbConfig.Host
	port := dbConfig.Port
	dbname := dbConfig.Name
	logsDbname := dbname + "_logs"

	// 1. Connect to Main DB
	MainDB = connectOrInitDB(user, pass, host, port, dbname)
	fmt.Println("Connected to Main PostgreSQL Database")

	// 2. Connect to Logs DB
	LogsDB = connectOrInitDB(user, pass, host, port, logsDbname)
	fmt.Println("Connected to Logs PostgreSQL Database")
}

func connectOrInitDB(user, pass, host, port, dbname string) *gorm.DB {
	// Standard Postgres DSN format
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, pass, dbname, port)

	// TranslateError surfaces portable sentinel errors (e.g. gorm.ErrDuplicatedKey)
	// instead of raw driver errors, so callers can reliably detect things like
	// unique constraint violations without parsing Postgres-specific error codes.
	gormConfig := &gorm.Config{TranslateError: true}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)

	if err != nil {
		// Postgres usually returns this exact error if the DB is missing
		if strings.Contains(err.Error(), "does not exist") || strings.Contains(err.Error(), "SQLSTATE 3D000") {
			createDatabase(user, pass, host, port, dbname)
			// Try connecting again after creation
			db, err = gorm.Open(postgres.Open(dsn), gormConfig)
		}
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL (%s): %v", dbname, err)
		}
	}
	return db
}

func createDatabase(user, pass, host, port, dbname string) {
	// In Postgres, you must connect to the default "postgres" database to issue a CREATE DATABASE command
	defaultDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/postgres?sslmode=disable", user, pass, host, port)

	sqlDB, err := sql.Open("pgx", defaultDSN)
	if err != nil {
		log.Fatalf("Failed to connect to default PostgreSQL server: %v", err)
	}
	defer sqlDB.Close()

	// CREATE DATABASE cannot run inside a transaction block in Postgres, so we use a raw Exec
	if _, err := sqlDB.Exec("CREATE DATABASE " + dbname); err != nil {
		log.Fatalf("Failed to create database %s: %v", dbname, err)
	}
	fmt.Printf("Successfully created database: %s\n", dbname)
}
