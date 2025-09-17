package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"web_backend/internal/app/ds"
	"web_backend/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Users{},
		&ds.Device{},
		&ds.Application{},
		&ds.ApplicationDevices{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}