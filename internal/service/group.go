package service

import (
	"context"
	"github/eventApp/internal/models"
)

type groupRep interface {
	CreateGroup(group *models.Group, ctx context.Context) (*models.Group, error)
	UpdateGroup(id int64, group *models.Group, ctx context.Context) (*models.Group, error)
	GetGroups(city, country string, ctx context.Context) ([]*models.Group, error)
}

type GroupService struct {
	groupRep groupRep
}

func NewGroupService(groupRep groupRep) *GroupService {
	return &GroupService{
		groupRep,
	}
}

type CreateGroupRequest struct {
	Name     string   `json:"name"`
	City     string   `json:"city"`
	Country  string   `json:"country"`
	KeyWords []string `json:"keyWords"`
}

type CreateGroupResponse struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	City     string   `json:"city"`
	Country  string   `json:"country"`
	KeyWords []string `json:"keyWords"`
}

func (s *GroupService) CreateGroup(cgr *CreateGroupRequest, ctx context.Context) (*CreateGroupResponse, error) {

	group := &models.Group{
		Name:     cgr.Name,
		City:     cgr.City,
		Country:  cgr.Country,
		KeyWords: cgr.KeyWords,
	}

	createdGroup, err := s.groupRep.CreateGroup(group, ctx)
	if err != nil {
		return nil, err
	}

	cgResp := &CreateGroupResponse{
		ID:       createdGroup.ID,
		Name:     createdGroup.Name,
		City:     createdGroup.City,
		Country:  createdGroup.Country,
		KeyWords: createdGroup.KeyWords,
	}

	return cgResp, nil

}

type GetGroupResponse struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	City     string   `json:"city"`
	Country  string   `json:"country"`
	KeyWords []string `json:"keyWords"`
}

func (s *GroupService) GetGroups(city, country string, ctx context.Context) ([]*GetGroupResponse, error) {

	groups, err := s.groupRep.GetGroups(city, country, ctx)
	if err != nil {
		return nil, err
	}

	groupsResp := make([]*GetGroupResponse, 0, len(groups))

	for _, g := range groups {
		groupsResp = append(groupsResp, &GetGroupResponse{
			ID:       g.ID,
			Name:     g.Name,
			City:     g.City,
			Country:  g.Country,
			KeyWords: g.KeyWords,
		})
	}

	return groupsResp, nil
}

type UpdateGroupRequest struct {
	Name     string   `json:"name"`
	City     string   `json:"city"`
	Country  string   `json:"country"`
	KeyWords []string `json:"keyWords"`
}

type UpdateGroupResponse struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	City     string   `json:"city"`
	Country  string   `json:"country"`
	KeyWords []string `json:"keyWords"`
}

func (s *GroupService) UpdateGroup(id int64, ugr *UpdateGroupRequest, ctx context.Context) (*UpdateGroupResponse, error) {

	group := &models.Group{
		Name:     ugr.Name,
		City:     ugr.City,
		Country:  ugr.Country,
		KeyWords: ugr.KeyWords,
	}

	updatedGroup, err := s.groupRep.UpdateGroup(id, group, ctx)
	if err != nil {
		return nil, err
	}

	ugResp := &UpdateGroupResponse{
		ID:       updatedGroup.ID,
		Name:     updatedGroup.Name,
		City:     updatedGroup.City,
		Country:  updatedGroup.Country,
		KeyWords: updatedGroup.KeyWords,
	}

	return ugResp, nil

}
