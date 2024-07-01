package config

import (
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

func NewSqlDb() (*gorm.DB, error) {
	addr := os.Getenv("DB_ADDRESS")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, addr, db)
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return gorm.Open(mysql.New(mysql.Config{Conn: sqlDB}), &gorm.Config{})
}
