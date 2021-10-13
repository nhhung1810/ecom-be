package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// https://github.com/jackc/pgx

var db *gorm.DB

func Connect() {
	//config
	dsn := "host=localhost user=postgres password=admin dbname=ecom port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	connection, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		println(err)
		panic("Could not connect to db")
	}

	db = connection
	print(db)

	connection.AutoMigrate(User{})
}

func GetDB() *gorm.DB {
	return db
}
