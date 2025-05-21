package repository

import (
	"context"
	"github/eventApp/internal/models"

	"github.com/uptrace/bun"
)

type GroupRepository struct {
	db *bun.DB
}

type Group struct {
	bun.BaseModel `bun:"table:groups,alias:u"`

	ID       int64  `bun:",pk,autoincrement,nullzero"`
	Name     string `bun:",unique"`
	City     string
	Country  string `bun:",unique"`
	KeyWords []string
	Events   []*Event `bun:"rel:has-many,join:id=group_id"`
	Users    []*User  `bun:"m2m:group_to_users,join:Group=User"`
}

func NewGroupRepository(db *bun.DB, ctx context.Context) (*GroupRepository, error) {
	usr := &GroupRepository{db}
	err := usr.createGroupTable(ctx)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *GroupRepository) createGroupTable(ctx context.Context) error {
	_, err := s.db.NewCreateTable().IfNotExists().Model((*Group)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *GroupRepository) CreateGroup(group *models.Group, ctx context.Context) (*models.Group, error) {

	g := &Group{
		Name:     group.Name,
		City:     group.City,
		Country:  group.Country,
		KeyWords: group.KeyWords,
	}

	createdGroup := &Group{}

	err := s.db.NewInsert().Model(g).Returning("*").Scan(ctx, createdGroup)
	if err != nil {
		return nil, err
	}

	cg := &models.Group{
		ID:       createdGroup.ID,
		Name:     createdGroup.Name,
		City:     createdGroup.City,
		Country:  createdGroup.Country,
		KeyWords: createdGroup.KeyWords,
	}

	return cg, nil
}

func (s *GroupRepository) UpdateGroup(id int64, group *models.Group, ctx context.Context) (*models.Group, error) {

	g := &Group{
		Name:     group.Name,
		City:     group.City,
		Country:  group.Country,
		KeyWords: group.KeyWords,
	}

	updatedGroup := &Group{}

	err := s.db.NewUpdate().Model(g).Where("id = ?", id).Returning("*").Scan(ctx, updatedGroup)
	if err != nil {
		return nil, err
	}

	ug := &models.Group{
		ID:       updatedGroup.ID,
		Name:     updatedGroup.Name,
		City:     updatedGroup.City,
		Country:  updatedGroup.Country,
		KeyWords: updatedGroup.KeyWords,
	}

	return ug, nil
}

func (s *GroupRepository) GetGroups(city, country string, ctx context.Context) ([]*models.Group, error) {
	var groups []Group

	err := s.db.NewSelect().Model(&groups).Where("city = ?", city).Where("country = ?", country).Scan(ctx)
	if err != nil {
		return nil, err
	}

	mgs := make([]*models.Group, 0, len(groups))

	for _, g := range groups {
		mgs = append(mgs, &models.Group{
			ID:       g.ID,
			Name:     g.Name,
			City:     g.City,
			Country:  g.Country,
			KeyWords: g.KeyWords,
		})
	}

	return mgs, nil
}
