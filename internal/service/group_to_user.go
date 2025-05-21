package service

import (
	"context"
	"github/eventApp/internal/models"
)

type groupToUserRep interface {
	AddUserToGroup(gtu *models.GroupToUser, ctx context.Context) (*models.GroupToUser, error)
	RemoveUserFromGroup(gtu *models.GroupToUser, ctx context.Context) error
}

type GroupToUserService struct {
	groupToUserRep groupToUserRep
}

func NewGroupToUserService(groupToUserRep groupToUserRep) *GroupToUserService {
	return &GroupToUserService{
		groupToUserRep,
	}
}

type AddUserToGroupRequest struct {
	GroupID int64 `json:"groupId"`
	UserID  int64 `json:"userId"`
}

type AddUserToGroupResponse struct {
	GroupID int64 `json:"groupId"`
	UserID  int64 `json:"userId"`
}

func (gtus *GroupToUserService) AddUserToGroup(autur *AddUserToGroupRequest, ctx context.Context) (*AddUserToGroupResponse, error) {

	gtu := &models.GroupToUser{
		GroupID: autur.GroupID,
		UserID:  autur.UserID,
	}

	addedUserToGroup, err := gtus.groupToUserRep.AddUserToGroup(gtu, ctx)
	if err != nil {
		return nil, err
	}

	autgResp := &AddUserToGroupResponse{
		GroupID: addedUserToGroup.GroupID,
		UserID:  addedUserToGroup.UserID,
	}

	return autgResp, nil

}

type RemoveUserFromGroupRequest struct {
	GroupID int64 `json:"groupId"`
	UserID  int64 `json:"userId"`
}

func (gtus *GroupToUserService) RemoveUserFromGroup(rufgr *RemoveUserFromGroupRequest, ctx context.Context) error {

	gtu := &models.GroupToUser{
		GroupID: rufgr.GroupID,
		UserID:  rufgr.UserID,
	}

	err := gtus.groupToUserRep.RemoveUserFromGroup(gtu, ctx)
	if err != nil {
		return err
	}

	return nil
}
