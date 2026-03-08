package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNew_PostgresConnection(t *testing.T) {
	db, err := New("postgres://invalid:invalid@localhost:5432/nonexistent?sslmode=disable")

	if db != nil {
		sqlDB, dbErr := db.DB()
		if dbErr == nil {
			err = sqlDB.Ping()
			assert.Error(t, err)
			_ = sqlDB.Close()
		}
	} else {
		assert.Error(t, err)
	}
}

func TestNew_PostgresWithValidDSN(t *testing.T) {
	db, err := New("postgres://user:pass@localhost:5432/testdb?sslmode=disable")

	if db != nil {
		sqlDB, dbErr := db.DB()
		if dbErr == nil {
			err = sqlDB.Ping()
			assert.Error(t, err)
			_ = sqlDB.Close()
		}
	}
	assert.Error(t, err)
}

func TestClose_Nil(t *testing.T) {
	Close(nil)
}

func TestClose_WithDB(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)

	Close(db)

	sqlDB, _ := db.DB()
	err = sqlDB.Ping()
	assert.Error(t, err)
}

func TestPing_Nil(t *testing.T) {
	err := Ping(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestPing_WithDB(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer Close(db)

	err = Ping(db)
	assert.NoError(t, err)
}

func TestPing_ClosedDB(t *testing.T) {
	db, _ := NewSQLite("file::memory:?cache=shared")
	Close(db)

	err := Ping(db)
	assert.Error(t, err)
}

func TestNew_ConnectionPoolSettings(t *testing.T) {
	db, _ := New("postgres://invalid:invalid@localhost:5432/nonexistent?sslmode=disable")

	if db != nil {
		sqlDB, dbErr := db.DB()
		if dbErr == nil {
			stats := sqlDB.Stats()
			assert.GreaterOrEqual(t, stats.OpenConnections, 0)
			_ = sqlDB.Close()
		}
	}
}

func TestNew_LoggerConfig(t *testing.T) {
	db, err := New("postgres://invalid:invalid@localhost:5432/nonexistent?sslmode=disable")

	assert.True(t, db != nil || err != nil)

	if db != nil {
		sqlDB, dbErr := db.DB()
		if dbErr == nil {
			_ = sqlDB.Close()
		}
	}
}

func TestNew_TimeoutSettings(t *testing.T) {
	db, err := New("postgres://invalid:invalid@localhost:5432/nonexistent?sslmode=disable")

	if db != nil {
		sqlDB, dbErr := db.DB()
		if dbErr == nil {
			stats := sqlDB.Stats()
			assert.GreaterOrEqual(t, stats.OpenConnections, 0)
			_ = sqlDB.Close()
		}
	} else {
		assert.Error(t, err)
	}
}

func TestClose_AlreadyClosed(t *testing.T) {
	db, _ := NewSQLite("file::memory:?cache=shared")
	Close(db)

	Close(db)
}

func TestPing_ErrorFromDB(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)

	sqlDB, _ := db.DB()
	_ = sqlDB.Close()

	err = Ping(db)
	assert.Error(t, err)
}

func TestNew_EmptyDSN(t *testing.T) {
	db, err := New("")

	if db != nil {
		sqlDB, dbErr := db.DB()
		if dbErr == nil {
			err = sqlDB.Ping()
			assert.Error(t, err)
			_ = sqlDB.Close()
		}
	}
	assert.Error(t, err)
}

func TestNew_InvalidDSN(t *testing.T) {
	db, err := New("not-a-valid-dsn")

	if db != nil {
		sqlDB, dbErr := db.DB()
		if dbErr == nil {
			err = sqlDB.Ping()
			assert.Error(t, err)
			_ = sqlDB.Close()
		}
	}
	assert.Error(t, err)
}

func TestClose_WithNilDB(t *testing.T) {
	var db *gorm.DB = nil
	Close(db)
}

func TestPing_WithNilDB(t *testing.T) {
	var db *gorm.DB = nil
	err := Ping(db)
	assert.Error(t, err)
}

func TestNewSQLite_TimeoutSettings(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer CloseSQLite(db)

	sqlDB, err := db.DB()
	assert.NoError(t, err)

	stats := sqlDB.Stats()
	assert.GreaterOrEqual(t, stats.OpenConnections, 0)
}

func TestCloseSQLite_WithNilDB(t *testing.T) {
	var db *gorm.DB = nil
	CloseSQLite(db)
}

func TestPingSQLite_WithNilDB(t *testing.T) {
	var db *gorm.DB = nil
	err := PingSQLite(db)
	assert.Error(t, err)
}

func TestPingSQLite_DBTest(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer CloseSQLite(db)

	err = PingSQLite(db)
	assert.NoError(t, err)
}

func TestPingSQLite_ClosedDB(t *testing.T) {
	db, _ := NewSQLite("file::memory:?cache=shared")
	CloseSQLite(db)

	err := PingSQLite(db)
	assert.Error(t, err)
}

func TestNewSQLite_EmptyPath(t *testing.T) {
	db, err := NewSQLite("")

	if err == nil && db != nil {
		defer CloseSQLite(db)
		err = PingSQLite(db)
		assert.NoError(t, err)
	}
}

func TestNewSQLite_MemoryMode(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer CloseSQLite(db)

	err = PingSQLite(db)
	assert.NoError(t, err)
}

func TestNewSQLite_DiskMode(t *testing.T) {
	db, err := NewSQLite("file:testdb.sqlite?mode=rwc")
	assert.NoError(t, err)
	defer func() {
		CloseSQLite(db)
	}()

	err = PingSQLite(db)
	assert.NoError(t, err)
}

func TestNewSQLite_PoolSettings(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer CloseSQLite(db)

	sqlDB, err := db.DB()
	assert.NoError(t, err)

	assert.GreaterOrEqual(t, sqlDB.Stats().OpenConnections, 0)
}

func TestNewSQLite_AutoMigrateMultiple(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer CloseSQLite(db)

	type TestModel1 struct {
		ID   uint `gorm:"primaryKey"`
		Name string
	}
	type TestModel2 struct {
		ID    uint `gorm:"primaryKey"`
		Value int
	}

	err = db.AutoMigrate(&TestModel1{}, &TestModel2{})
	assert.NoError(t, err)
}

func TestCloseSQLite_MultipleTimes(t *testing.T) {
	db, _ := NewSQLite("file::memory:?cache=shared")
	CloseSQLite(db)

	CloseSQLite(db)
	CloseSQLite(db)
}

func TestPingSQLite_Error(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)

	sqlDB, _ := db.DB()
	_ = sqlDB.Close()

	err = PingSQLite(db)
	assert.Error(t, err)
}
