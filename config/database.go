package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Sqllogger struct {
	logger.Interface
}

func (l Sqllogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v\n================================\n", sql)
}

func SetupDatabaseConnection() *gorm.DB {
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		panic("Failed to load env file")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dbCont := fmt.Sprintf("%s:%s@tcp(%s:3307)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dbCont), &gorm.Config{
		Logger: &Sqllogger{},
	})
	if err != nil {
		panic("Failed to  connect to database")
	}

	return db
}
