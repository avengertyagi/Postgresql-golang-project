package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	permissionmodel "github.com/akshit_tyagi/postgresql_project/internal/models/permission"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

type Database struct {
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
}

func LoadDatabase() Database {
	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}
	maxOpen, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	if maxOpen == 0 {
		maxOpen = 25
	}
	maxIdle, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	if maxIdle == 0 {
		maxIdle = 25
	}
	connLife, _ := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))
	if connLife == 0 {
		connLife = 5 * time.Minute
	}
	return Database{
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBUser:            os.Getenv("DB_USERNAME"),
		DBPassword:        os.Getenv("DB_PASSWORD"),
		DBName:            os.Getenv("DB_DATABASE"),
		DBSSLMode:         sslMode,
		DBMaxOpenConns:    maxOpen,
		DBMaxIdleConns:    maxIdle,
		DBConnMaxLifetime: connLife,
	}
}

func (c Database) GetDBConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode,
	)
}

func (c Database) GetDBConnectionURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}

func InitializeDatabase() error {
	dbConfig := LoadDatabase()
	if dbConfig.DBHost == "" || dbConfig.DBName == "" {
		return errors.New("database: DB_HOST or DB_DATABASE env vars are not set")
	}
	var gormLogger logger.Interface
	if os.Getenv("APP_ENV") == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Warn)
	}
	db, err := gorm.Open(postgres.Open(dbConfig.GetDBConnectionString()),
		&gorm.Config{Logger: gormLogger},
	)
	if err != nil {
		return err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("database: get sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(dbConfig.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(dbConfig.DBConnMaxLifetime)
	db.AutoMigrate(&permissionmodel.Permission)
	DB = db
	return nil
}

func Close() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
