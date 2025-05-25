package service

import (
	"context"
	"fmt"
	"github/eventApp/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type userRep interface {
	CreateUser(user *models.User, ctx context.Context) (*models.User, error)
	UpdateUser(id int64, user *models.User, ctx context.Context) (*models.User, error)
	GetUser(id int64, ctx context.Context) (*models.User, error)
}

type UserService struct {
	userRep userRep
}

func NewUserService(userRep userRep) *UserService {
	return &UserService{
		userRep,
	}
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"userName"`
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (s *UserService) CreateUser(cur *CreateUserRequest, ctx context.Context) (*CreateUserResponse, error) {

	if cur.Password == "" {
		return nil, fmt.Errorf("password has to have a non-zero length")
	}

	hashedPassword, err := hashPassword(cur.Password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password")
	}

	user := &models.User{
		Name:     cur.Name,
		Email:    cur.Email,
		UserName: cur.UserName,
		Password: hashedPassword,
	}

	createdUser, err := s.userRep.CreateUser(user, ctx)
	if err != nil {
		return nil, err
	}

	cuResp := &CreateUserResponse{
		ID:       createdUser.ID,
		Name:     createdUser.Name,
		Email:    createdUser.Email,
		UserName: createdUser.UserName,
	}

	return cuResp, nil

}

type GetUserResponse struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	UserName string `json:"userName"`
}

func (s *UserService) GetUser(id int64, ctx context.Context) (*GetUserResponse, error) {

	user, err := s.userRep.GetUser(id, ctx)
	if err != nil {
		return nil, err
	}

	guResp := &GetUserResponse{
		Name:     user.Name,
		Email:    user.Email,
		UserName: user.UserName,
	}

	return guResp, nil
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *UserService) UpdateUser(id int64, uur *UpdateUserRequest, ctx context.Context) (*UpdateUserResponse, error) {

	user := &models.User{
		Name:  uur.Name,
		Email: uur.Email,
	}

	updatedUser, err := s.userRep.UpdateUser(id, user, ctx)
	if err != nil {
		return nil, err
	}

	uuResp := &UpdateUserResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
	}

	return uuResp, nil

}
