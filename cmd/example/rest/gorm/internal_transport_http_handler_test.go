package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserDTO), args.Error(1)
}

func (m *MockUserService) Get(ctx context.Context, id uint) (*UserDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserDTO), args.Error(1)
}

func (m *MockUserService) List(ctx context.Context, page, pageSize int) ([]UserDTO, int64, error) {
	args := m.Called(ctx, page, pageSize)
	return args.Get(0).([]UserDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) Update(ctx context.Context, id uint, req UpdateUserRequest) (*UserDTO, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*UserDTO), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupHandlerTest() (*gin.Engine, *MockUserService, *Handler) {
	gin.SetMode(gin.TestMode)
	userSvc := &MockUserService{}
	handler := NewHandler(userSvc)
	_ = handler.Router("file::memory:?cache=shared")
	return handler.router, userSvc, handler
}

func TestNewHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userSvc := &MockUserService{}
	handler := NewHandler(userSvc)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.router)
	assert.NotNil(t, handler.validate)
	assert.Equal(t, userSvc, handler.userSvc)
}

func TestHandler_Router(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userSvc := &MockUserService{}
	handler := NewHandler(userSvc)

	router := handler.Router("file::memory:?cache=shared")

	assert.NotNil(t, router)

	routes := router.Routes()
	assert.GreaterOrEqual(t, len(routes), 5) // At least 5 user routes
}

func TestHandler_createUser_Success(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	req := CreateUserRequest{
		Name:     "John Doe",
		Email:    "john@example.com",
		Password: "securepassword123",
	}

	userSvc.On("Create", mock.Anything, req).Return(&UserDTO{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/users", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
	assert.Contains(t, w.Body.String(), "john@example.com")
	userSvc.AssertExpectations(t)
}

func TestHandler_createUser_AlreadyExists(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	req := CreateUserRequest{
		Name:     "John Doe",
		Email:    "existing@example.com",
		Password: "securepassword123",
	}

	userSvc.On("Create", mock.Anything, req).Return((*UserDTO)(nil), ErrAlreadyExists)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/users", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "already exists")
	userSvc.AssertExpectations(t)
}

func TestHandler_createUser_InvalidJSON(t *testing.T) {
	router, _, _ := setupHandlerTest()

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/users", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_createUser_RequestTooLarge(t *testing.T) {
	router, _, _ := setupHandlerTest()

	largePassword := strings.Repeat("a", 1024*1024)
	payload := fmt.Sprintf(`{"name":"John Doe","email":"john@example.com","password":"%s"}`, largePassword)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/users", bytes.NewReader([]byte(payload)))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	assert.Contains(t, w.Body.String(), "request body too large")
}

func TestHandler_createUser_ValidationFailed(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	req := CreateUserRequest{
		Name:     "J", // too short
		Email:    "invalid",
		Password: "short", // too short
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/api/users", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	userSvc.AssertNotCalled(t, "Create")
}

func TestHandler_listUsers_Success(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	dtos := []UserDTO{
		{ID: 1, Name: "User 1", Email: "user1@example.com"},
		{ID: 2, Name: "User 2", Email: "user2@example.com"},
	}

	userSvc.On("List", mock.Anything, 1, 20).Return(dtos, int64(100), nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/users", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "items")
	assert.Contains(t, response, "meta")
	userSvc.AssertExpectations(t)
}

func TestHandler_listUsers_WithPagination(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	dtos := []UserDTO{
		{ID: 1, Name: "User 1", Email: "user1@example.com"},
	}

	userSvc.On("List", mock.Anything, 2, 10).Return(dtos, int64(50), nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/users?page=2&size=10", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	userSvc.AssertExpectations(t)
}

func TestHandler_listUsers_InvalidPage(t *testing.T) {
	router, _, _ := setupHandlerTest()

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/users?page=invalid", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid page")
}

func TestHandler_listUsers_InvalidSize(t *testing.T) {
	router, _, _ := setupHandlerTest()

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/users?size=invalid", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid size")
}

func TestHandler_getUser_Success(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	userSvc.On("Get", mock.Anything, uint(1)).Return(&UserDTO{
		ID:    1,
		Name:  "John Doe",
		Email: "john@example.com",
	}, nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/users/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "John Doe")
	userSvc.AssertExpectations(t)
}

func TestHandler_getUser_NotFound(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	userSvc.On("Get", mock.Anything, uint(999)).Return((*UserDTO)(nil), ErrNotFound)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/users/999", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
	userSvc.AssertExpectations(t)
}

func TestHandler_getUser_InvalidID(t *testing.T) {
	router, _, _ := setupHandlerTest()

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("GET", "/api/users/invalid", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid id")
}

func TestHandler_updateUser_Success(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	req := UpdateUserRequest{
		Name: new("Updated Name"),
	}

	userSvc.On("Update", mock.Anything, uint(1), req).Return(&UserDTO{
		ID:    1,
		Name:  "Updated Name",
		Email: "john@example.com",
	}, nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/users/1", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Updated Name")
	userSvc.AssertExpectations(t)
}

func TestHandler_updateUser_InvalidID(t *testing.T) {
	router, _, _ := setupHandlerTest()

	req := UpdateUserRequest{
		Name: new("Updated Name"),
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/api/users/invalid", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_deleteUser_Success(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	userSvc.On("Delete", mock.Anything, uint(1)).Return(nil)

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/users/1", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNoContent, w.Code)
	userSvc.AssertExpectations(t)
}

func TestHandler_deleteUser_Error(t *testing.T) {
	router, userSvc, _ := setupHandlerTest()

	userSvc.On("Delete", mock.Anything, uint(999)).Return(error(nil))

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/users/999", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNoContent, w.Code)
	userSvc.AssertExpectations(t)
}

func TestHandler_deleteUser_InvalidID(t *testing.T) {
	router, _, _ := setupHandlerTest()

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("DELETE", "/api/users/invalid", nil)

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_parseUintParam(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userSvc := &MockUserService{}
	handler := NewHandler(userSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "123"}}

	id, ok := handler.parseUintParam(c, "id")

	assert.True(t, ok)
	assert.Equal(t, uint(123), id)
}

func TestHandler_parseUintParam_Invalid(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userSvc := &MockUserService{}
	handler := NewHandler(userSvc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	id, ok := handler.parseUintParam(c, "id")

	assert.False(t, ok)
	assert.Equal(t, uint(0), id)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandler_bindAndValidate_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userSvc := &MockUserService{}
	handler := NewHandler(userSvc)

	assert.NotNil(t, handler)
	assert.NotNil(t, handler.validate)
}

func TestJSONError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	JSONError(c, http.StatusBadRequest, "test error")

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "test error")
	assert.Contains(t, w.Body.String(), "error")
}

func TestNotFoundHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	NotFoundHandler(c)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "resource not found")
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userSvc := &MockUserService{}
	handler := NewHandler(userSvc)
	_ = handler.Router("file::memory:?cache=shared")

	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PATCH", "/api/users/1", nil)

	handler.router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
