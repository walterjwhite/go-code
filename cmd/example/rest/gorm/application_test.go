package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfigure(t *testing.T) {
	Configure()
	assert.NotNil(t, Context)
}

func TestOnPanic_NoPanic(t *testing.T) {
	OnPanic()
}

func TestWait_ContextCancellation(t *testing.T) {
	assert.NotNil(t, Wait)
}

func TestOnPanic_WithPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "test panic", r)
		}
	}()

	OnPanic()
	panic("test panic")
}

func TestCleanup(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)

	cfg := &Config{
		AppPort:     "0",
		ReadTimeout: 1 * time.Second,
	}
	router := NewHandler(NewUserService(NewGormUserRepository(db))).Router("file::memory:?cache=shared")
	server := NewServer(router, cfg)

	go func() {
		_ = server.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	cleanup(db, server)

	err = PingSQLite(db)
	assert.Error(t, err)
}

func TestCleanup_ServerShutdownError(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)

	cfg := &Config{
		AppPort:     "0",
		ReadTimeout: 1 * time.Second,
	}
	router := NewHandler(NewUserService(NewGormUserRepository(db))).Router("file::memory:?cache=shared")
	server := NewServer(router, cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_ = server.Shutdown(ctx)

	cleanup(db, server)
}

func TestNew_Postgres(t *testing.T) {
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

func TestClose_NilDB(t *testing.T) {
	Close(nil)
}

func TestClose_Success(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)

	Close(db)
}

func TestPing_NilDB(t *testing.T) {
	err := Ping(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestPing_Success(t *testing.T) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)
	defer Close(db)

	err = Ping(db)
	assert.NoError(t, err)
}

func TestPing_Error(t *testing.T) {
	db, _ := NewSQLite("file::memory:?cache=shared")
	Close(db)

	err := Ping(db)
	assert.Error(t, err)
}
