package main

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewSQLite(databasePath string) (*gorm.DB, error) {
	gormLogger := logger.Default.LogMode(logger.Silent)

	dialector := sqlite.Open(databasePath)
	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(4)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}

func CloseSQLite(db *gorm.DB) {
	if db == nil {
		return
	}
	sqlDB, err := db.DB()
	if err != nil {
		return
	}
	_ = sqlDB.Close()
}

func PingSQLite(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
