package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/tangoctrung/golang_api_v2/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	errEnv := godotenv.Load()

	if errEnv != nil {
		panic("Failed to load env file")
	}

	DB_USER := os.Getenv("DB_USER")
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	DB_HOST := os.Getenv("DB_HOST")
	DB_NAME := os.Getenv("DB_NAME")

	url_DB := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", DB_USER, DB_PASSWORD, DB_HOST, DB_NAME)
	db, err := gorm.Open(mysql.Open(url_DB), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	fmt.Sprintln("Connect database success!")
	db.AutoMigrate(&entity.Book{}, &entity.User{})

	return db

}

func CloseConnectDatabase(db *gorm.DB) {
	dbSQL, err := db.DB()

	if err != nil {
		panic("Failed to close connect database")
	}

	dbSQL.Close()
}
