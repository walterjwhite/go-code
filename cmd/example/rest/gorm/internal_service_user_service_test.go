package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, u *User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context, offset, limit int) ([]User, error) {
	args := m.Called(ctx, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, u *User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) Count(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func TestNewUserService(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)

	assert.NotNil(t, svc)
	assert.IsType(t, &userService{}, svc)
}

func TestUserService_Create_Success(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	req := CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepassword123",
	}

	repo.On("GetByEmail", ctx, "john@example.com").Return((*User)(nil), ErrNotFound)
	repo.On("Create", ctx, mock.MatchedBy(func(u *User) bool {
		return u.Name == "John Doe" && u.Email == "john@example.com" && u.PasswordHash != ""
	})).Run(func(args mock.Arguments) {
		u := args.Get(1).(*User)
		u.ID = 1
	}).Return(nil)

	dto, err := svc.Create(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Equal(t, "John Doe", dto.Name)
	assert.Equal(t, "john@example.com", dto.Email)
	assert.NotZero(t, dto.ID)
	repo.AssertExpectations(t)
}

func TestUserService_Create_AlreadyExists(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	req := CreateUserRequest{
		Name:     "John Doe",
		Email:    "existing@example.com",
		Password: "securepassword123",
	}

	existingUser := &User{
		ID:           1,
		Name:         "Existing",
		Email:        "existing@example.com",
		PasswordHash: "hashed",
	}

	repo.On("GetByEmail", ctx, "existing@example.com").Return(existingUser, nil)

	dto, err := svc.Create(ctx, req)

	assert.ErrorIs(t, err, ErrAlreadyExists)
	assert.Nil(t, dto)
	repo.AssertExpectations(t)
}

func TestUserService_Create_RepositoryError(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	req := CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepassword123",
	}

	repo.On("GetByEmail", ctx, "john@example.com").Return((*User)(nil), ErrNotFound)
	repo.On("Create", ctx, mock.Anything).Return(errors.New("db error"))

	dto, err := svc.Create(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	assert.Nil(t, dto)
	repo.AssertExpectations(t)
}

func TestUserService_Get_Success(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	existingUser := &User{
		ID:           1,
		Name:         "John Doe",
		Email:        "john@example.com",
		PasswordHash: "hashed",
	}

	repo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)

	dto, err := svc.Get(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Equal(t, uint(1), dto.ID)
	assert.Equal(t, "John Doe", dto.Name)
	assert.Equal(t, "john@example.com", dto.Email)
	repo.AssertExpectations(t)
}

func TestUserService_Get_NotFound(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, uint(999)).Return((*User)(nil), ErrNotFound)

	dto, err := svc.Get(ctx, 999)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, dto)
	repo.AssertExpectations(t)
}

func TestUserService_List_Success(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	users := []User{
		{ID: 1, Name: "User 1", Email: "user1@example.com", PasswordHash: "hashed"},
		{ID: 2, Name: "User 2", Email: "user2@example.com", PasswordHash: "hashed"},
	}

	repo.On("List", ctx, 0, 20).Return(users, nil)
	repo.On("Count", ctx).Return(int64(100), nil)

	dtos, total, err := svc.List(ctx, 1, 20)

	assert.NoError(t, err)
	assert.Len(t, dtos, 2)
	assert.Equal(t, int64(100), total)
	assert.Equal(t, "User 1", dtos[0].Name)
	assert.Equal(t, "User 2", dtos[1].Name)
	repo.AssertExpectations(t)
}

func TestUserService_List_PageValidation(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	users := []User{
		{ID: 1, Name: "User 1", Email: "user1@example.com", PasswordHash: "hashed"},
	}

	repo.On("List", ctx, 0, 20).Return(users, nil)
	repo.On("Count", ctx).Return(int64(1), nil)

	dtos, total, err := svc.List(ctx, 0, 20)

	assert.NoError(t, err)
	assert.Len(t, dtos, 1)
	assert.Equal(t, int64(1), total)
	repo.AssertExpectations(t)
}

func TestUserService_List_PageSizeValidation(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	users := []User{
		{ID: 1, Name: "User 1", Email: "user1@example.com", PasswordHash: "hashed"},
	}

	repo.On("List", ctx, 0, 20).Return(users, nil)
	repo.On("Count", ctx).Return(int64(1), nil)

	dtos, _, err := svc.List(ctx, 1, 0)

	assert.NoError(t, err)
	assert.Len(t, dtos, 1)
	repo.AssertExpectations(t)
}

func TestUserService_List_PageSizeTooLarge(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	users := []User{
		{ID: 1, Name: "User 1", Email: "user1@example.com", PasswordHash: "hashed"},
	}

	repo.On("List", ctx, 0, 20).Return(users, nil)
	repo.On("Count", ctx).Return(int64(1), nil)

	dtos, _, err := svc.List(ctx, 1, 150)

	assert.NoError(t, err)
	assert.Len(t, dtos, 1)
	repo.AssertExpectations(t)
}

func TestUserService_List_RepositoryError(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	repo.On("List", ctx, 0, 20).Return([]User(nil), errors.New("list error"))

	dtos, total, err := svc.List(ctx, 1, 20)

	assert.Error(t, err)
	assert.Equal(t, "list error", err.Error())
	assert.Nil(t, dtos)
	assert.Equal(t, int64(0), total)
	repo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	existingUser := &User{
		ID:           1,
		Name:         "Original Name",
		Email:        "original@example.com",
		PasswordHash: "hashed",
	}

	repo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
	repo.On("Update", ctx, existingUser).Return(nil)

	req := UpdateUserRequest{
		Name: new("Updated Name"),
	}

	dto, err := svc.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Equal(t, "Updated Name", dto.Name)
	assert.Equal(t, "original@example.com", dto.Email)
	repo.AssertExpectations(t)
}

func TestUserService_Update_EmailOnly(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	existingUser := &User{
		ID:           1,
		Name:         "Original Name",
		Email:        "original@example.com",
		PasswordHash: "hashed",
	}

	repo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
	repo.On("Update", ctx, existingUser).Return(nil)

	req := UpdateUserRequest{
		Email: new("updated@example.com"),
	}

	dto, err := svc.Update(ctx, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, dto)
	assert.Equal(t, "Original Name", dto.Name)
	assert.Equal(t, "updated@example.com", dto.Email)
	repo.AssertExpectations(t)
}

func TestUserService_Update_NotFound(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	repo.On("GetByID", ctx, uint(999)).Return((*User)(nil), ErrNotFound)

	req := UpdateUserRequest{
		Name: new("New Name"),
	}

	dto, err := svc.Update(ctx, 999, req)

	assert.ErrorIs(t, err, ErrNotFound)
	assert.Nil(t, dto)
	repo.AssertExpectations(t)
}

func TestUserService_Update_RepositoryError(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	existingUser := &User{
		ID:           1,
		Name:         "Original",
		Email:        "original@example.com",
		PasswordHash: "hashed",
	}

	repo.On("GetByID", ctx, uint(1)).Return(existingUser, nil)
	repo.On("Update", ctx, existingUser).Return(errors.New("update error"))

	req := UpdateUserRequest{
		Name: new("New Name"),
	}

	dto, err := svc.Update(ctx, 1, req)

	assert.Error(t, err)
	assert.Equal(t, "update error", err.Error())
	assert.Nil(t, dto)
	repo.AssertExpectations(t)
}

func TestUserService_Delete_Success(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	repo.On("Delete", ctx, uint(1)).Return(nil)

	err := svc.Delete(ctx, 1)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestUserService_Delete_Error(t *testing.T) {
	repo := &MockUserRepository{}
	svc := NewUserService(repo)
	ctx := context.Background()

	repo.On("Delete", ctx, uint(999)).Return(errors.New("delete error"))

	err := svc.Delete(ctx, 999)

	assert.Error(t, err)
	assert.Equal(t, "delete error", err.Error())
	repo.AssertExpectations(t)
}

func TestUserToDTO_NilUser(t *testing.T) {
	dto := userToDTO(nil)
	assert.Nil(t, dto)
}

func TestUserToDTO_ValidUser(t *testing.T) {
	user := &User{
		ID:           1,
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "secret",
	}

	dto := userToDTO(user)

	assert.NotNil(t, dto)
	assert.Equal(t, uint(1), dto.ID)
	assert.Equal(t, "Test User", dto.Name)
	assert.Equal(t, "test@example.com", dto.Email)
}

func TestHashPassword(t *testing.T) {
	password := "securepassword123"

	hash, err := hashPassword(password)

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)
}

func TestCreateUserRequest_Validation(t *testing.T) {
	req := CreateUserRequest{
		Name:     "Jo", // min=2
		Email:    "invalid-email",
		Password: "short", // min=8
	}

	assert.Equal(t, "Jo", req.Name)
	assert.Equal(t, "invalid-email", req.Email)
	assert.Equal(t, "short", req.Password)
}

func TestUpdateUserRequest_NilFields(t *testing.T) {
	req := UpdateUserRequest{
		Name:  nil,
		Email: nil,
	}

	assert.Nil(t, req.Name)
	assert.Nil(t, req.Email)
}
