package config

import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    var err error
    dsn := "user=postgres password=1144 dbname=kuplyu_app port=5432 sslmode=disable TimeZone=Asia/Shanghai"
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        panic("Failed to connect to the database")
    }
}
