package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:@tcp(localhost:3306)/lms_usti?charset=utf8mb4&parseTime=True&loc=Local"
	database, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn, DefaultStringSize: 255}), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err.Error())
	}
	database.AutoMigrate(&User{}, &VerificationToken{}, &Classroom{}, &Announcement{}, &Material{}, &MaterialFile{}, &MaterialLink{})
	DB = database
}
