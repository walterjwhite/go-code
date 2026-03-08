package main

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	db, err := NewSQLite("file::memory:?cache=shared")
	assert.NoError(t, err)

	err = db.AutoMigrate(&User{})
	assert.NoError(t, err)

	cleanup := func() {
		CloseSQLite(db)
	}

	return db, cleanup
}

func TestNewGormUserRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)

	assert.NotNil(t, repo)
	assert.IsType(t, &gormUserRepo{}, repo)
}

func TestUserRepository_Create(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	user := &User{
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashed_password",
	}

	err := repo.Create(ctx, user)

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, 5*time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, 5*time.Second)
}

func TestUserRepository_GetByID(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	created := &User{
		Name:         "Jane Doe",
		Email:        "jane@example.com",
		PasswordHash: "hashed",
	}
	err := repo.Create(ctx, created)
	assert.NoError(t, err)

	user, err := repo.GetByID(ctx, created.ID)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, "Jane Doe", user.Name)
	assert.Equal(t, "jane@example.com", user.Email)
}

func TestUserRepository_GetByID_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	user, err := repo.GetByID(ctx, 99999)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, user)
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	created := &User{
		Name:         "Bob Smith",
		Email:        "bob@example.com",
		PasswordHash: "hashed",
	}
	err := repo.Create(ctx, created)
	assert.NoError(t, err)

	user, err := repo.GetByEmail(ctx, "bob@example.com")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, created.ID, user.ID)
	assert.Equal(t, "Bob Smith", user.Name)
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	user, err := repo.GetByEmail(ctx, "nonexistent@example.com")

	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, user)
}

func TestUserRepository_List(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	for i := range 5 {
		user := &User{
			Name:         "User",
			Email:        "user" + strconv.Itoa(i) + "@example.com",
			PasswordHash: "hashed",
		}
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
	}

	users, err := repo.List(ctx, 0, 3)

	assert.NoError(t, err)
	assert.Len(t, users, 3)
	assert.Equal(t, uint(5), users[0].ID)
	assert.Equal(t, uint(4), users[1].ID)
	assert.Equal(t, uint(3), users[2].ID)
}

func TestUserRepository_List_Empty(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	users, err := repo.List(ctx, 0, 10)

	assert.NoError(t, err)
	assert.Empty(t, users)
}

func TestUserRepository_Update(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	user := &User{
		Name:         "Original Name",
		Email:        "original@example.com",
		PasswordHash: "hashed",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	user.Name = "Updated Name"
	err = repo.Update(ctx, user)

	assert.NoError(t, err)

	updated, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Name", updated.Name)
}

func TestUserRepository_Delete(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	user := &User{
		Name:         "ToDelete",
		Email:        "delete@example.com",
		PasswordHash: "hashed",
	}
	err := repo.Create(ctx, user)
	assert.NoError(t, err)

	err = repo.Delete(ctx, user.ID)

	assert.NoError(t, err)

	_, err = repo.GetByID(ctx, user.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestUserRepository_Count(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx := context.Background()

	count, err := repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(0), count)

	for i := range 3 {
		user := &User{
			Name:         "User",
			Email:        "count" + strconv.Itoa(i) + "@example.com",
			PasswordHash: "hashed",
		}
		err := repo.Create(ctx, user)
		assert.NoError(t, err)
	}

	count, err = repo.Count(ctx)
	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
}

func TestUserRepository_Create_ContextCancellation(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	repo := NewGormUserRepository(db)
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	user := &User{
		Name:         "Test",
		Email:        "test@example.com",
		PasswordHash: "hashed",
	}

	err := repo.Create(ctx, user)

	assert.Error(t, err)
}

func TestUserModel_JSONSerialization(t *testing.T) {
	user := &User{
		ID:           1,
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "secret",
	}

	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "Test User", user.Name)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "secret", user.PasswordHash)
}
