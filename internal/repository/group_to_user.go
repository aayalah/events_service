package repository

import (
	"context"
	"github/eventApp/internal/models"

	"github.com/uptrace/bun"
)

type GroupToUserRepository struct {
	db *bun.DB
}

type GroupToUser struct {
	GroupID int64  `bun:",pk"`
	Group   *Group `bun:"rel:belongs-to,join:group_id=id"`
	UserID  int64  `bun:",pk"`
	User    *User  `bun:"rel:belongs-to,join:user_id=id"`
}

func NewGroupToUserRepository(db *bun.DB, ctx context.Context) (*GroupToUserRepository, error) {
	gtur := &GroupToUserRepository{db}
	err := gtur.createGroupToUserTable(ctx)
	if err != nil {
		return nil, err
	}
	db.RegisterModel((*GroupToUser)(nil))

	return gtur, nil
}

func (gtur *GroupToUserRepository) createGroupToUserTable(ctx context.Context) error {
	_, err := gtur.db.NewCreateTable().IfNotExists().Model((*GroupToUser)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (gtur *GroupToUserRepository) AddUserToGroup(gtu *models.GroupToUser, ctx context.Context) (*models.GroupToUser, error) {
	gu := &GroupToUser{
		GroupID: gtu.GroupID,
		UserID:  gtu.UserID,
	}

	createdGroupToUser := &GroupToUser{}

	err := gtur.db.NewInsert().Model(gu).Returning("*").Scan(ctx, createdGroupToUser)
	if err != nil {
		return nil, err
	}

	gum := &models.GroupToUser{
		GroupID: createdGroupToUser.GroupID,
		UserID:  createdGroupToUser.UserID,
	}

	return gum, nil
}

func (gtur *GroupToUserRepository) RemoveUserFromGroup(gtu *models.GroupToUser, ctx context.Context) error {
	_, err := gtur.db.NewDelete().Model(&GroupToUser{}).Where("group_id = ?", gtu.GroupID).Where("user_id = ?", gtu.UserID).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
