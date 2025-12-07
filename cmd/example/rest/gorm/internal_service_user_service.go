package main

import (
	"context"
	"errors"


	"golang.org/x/crypto/bcrypt"
)

var ErrAlreadyExists = errors.New("already exists")

type UserDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name" validate:"omitempty,min=2,max=100"`
	Email *string `json:"email" validate:"omitempty,email"`
}

type UserService interface {
	Create(ctx context.Context, req CreateUserRequest) (*UserDTO, error)
	Get(ctx context.Context, id uint) (*UserDTO, error)
	List(ctx context.Context, page, pageSize int) ([]UserDTO, int64, error)
	Update(ctx context.Context, id uint, req UpdateUserRequest) (*UserDTO, error)
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) Create(ctx context.Context, req CreateUserRequest) (*UserDTO, error) {
	if _, err := s.repo.GetByEmail(ctx, req.Email); err == nil {
		return nil, ErrAlreadyExists
	} else if err != ErrNotFound {
		if err != nil {
			return nil, err
		}
	}

	hash, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hash,
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return userToDTO(u), nil
}

func (s *userService) Get(ctx context.Context, id uint) (*UserDTO, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return userToDTO(u), nil
}

func (s *userService) List(ctx context.Context, page, pageSize int) ([]UserDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	users, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.repo.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	dtos := make([]UserDTO, len(users))
	for i, u := range users {
		dtos[i] = *userToDTO(&u)
	}
	return dtos, total, nil
}

func (s *userService) Update(ctx context.Context, id uint, req UpdateUserRequest) (*UserDTO, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if err := s.repo.Update(ctx, u); err != nil {
		return nil, err
	}
	return userToDTO(u), nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func hashPassword(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	return string(b), err
}

func userToDTO(u *User) *UserDTO {
	if u == nil {
		return nil
	}
	return &UserDTO{ID: u.ID, Name: u.Name, Email: u.Email}
}
