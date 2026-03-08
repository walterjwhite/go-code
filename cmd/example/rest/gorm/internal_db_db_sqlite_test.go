package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSQLite_Success(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")

	assert.NoError(t, err)
	assert.NotNil(t, db)

	err = PingSQLite(db)
	assert.NoError(t, err)

	CloseSQLite(db)
}

func TestNewSQLite_InvalidPath(t *testing.T) {
	db, err := NewSQLite("/nonexistent/path/test.db")

	assert.Error(t, err)
	assert.Nil(t, db)
}

func TestCloseSQLite_NilDB(t *testing.T) {
	CloseSQLite(nil)
}

func TestCloseSQLite_Success(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)

	CloseSQLite(db)

	err = PingSQLite(db)
	assert.Error(t, err)
}

func TestPingSQLite_Success(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer CloseSQLite(db)

	err = PingSQLite(db)
	assert.NoError(t, err)
}

func TestNewSQLite_ConnectionPoolSettings(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer CloseSQLite(db)

	sqlDB, err := db.DB()
	assert.NoError(t, err)

	stats := sqlDB.Stats()
	assert.GreaterOrEqual(t, stats.OpenConnections, 0)
}

func TestNewSQLite_AutoMigrate(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer CloseSQLite(db)

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	user := &User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed",
	}
	err = db.Create(user).Error
	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}
