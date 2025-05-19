package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() (*gorm.DB, error) {
	dsn := "host=ep-odd-mud-aau4cvww-pooler.westus3.azure.neon.tech user=neondb_owner password=npg_H20AghUkNFTY dbname=neondb port=5432 sslmode=require TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Submission{})
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

func Disconnect(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.Close()
}
