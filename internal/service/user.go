package service

import (
	"context"
	"github/eventApp/internal/models"
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
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *UserService) CreateUser(cur *CreateUserRequest, ctx context.Context) (*CreateUserResponse, error) {

	user := &models.User{
		Name:  cur.Name,
		Email: cur.Email,
	}

	createdUser, err := s.userRep.CreateUser(user, ctx)
	if err != nil {
		return nil, err
	}

	cuResp := &CreateUserResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	}

	return cuResp, nil

}

type GetUserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *UserService) GetUser(id int64, ctx context.Context) (*GetUserResponse, error) {

	user, err := s.userRep.GetUser(id, ctx)
	if err != nil {
		return nil, err
	}

	guResp := &GetUserResponse{
		Name:  user.Name,
		Email: user.Email,
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
