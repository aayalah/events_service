package repository

import (
	"context"
	"github/eventApp/internal/models"

	"github.com/uptrace/bun"
)

type UserRepository struct {
	db *bun.DB
}

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID       int64  `bun:",pk,autoincrement,nullzero"`
	Name     string `bun:",notnull,unique"`
	Email    string `bun:",notnull,unique"`
	UserName string `bun:",notnull,unique"`
	Password string `bun:",notnull"`
}

func NewUserRepository(db *bun.DB, ctx context.Context) (*UserRepository, error) {
	usr := &UserRepository{db}
	err := usr.createUserTable(ctx)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

func (s *UserRepository) createUserTable(ctx context.Context) error {
	_, err := s.db.NewCreateTable().IfNotExists().Model((*User)(nil)).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserRepository) CreateUser(user *models.User, ctx context.Context) (*models.User, error) {

	us := &User{
		Name:     user.Name,
		Email:    user.Email,
		UserName: user.UserName,
		Password: user.Password,
	}

	createdUser := &User{}

	err := s.db.NewInsert().Model(us).Returning("*").Scan(ctx, createdUser)
	if err != nil {
		return nil, err
	}

	ud := &models.User{
		ID:       createdUser.ID,
		Name:     createdUser.Name,
		Email:    createdUser.Email,
		UserName: createdUser.UserName,
	}

	return ud, nil
}

func (s *UserRepository) UpdateUser(id int64, user *models.User, ctx context.Context) (*models.User, error) {

	us := &User{
		Name:  user.Name,
		Email: user.Email,
	}

	updatedUser := &User{}

	err := s.db.NewUpdate().Model(us).Where("id = ?", id).Returning("*").Scan(ctx, updatedUser)
	if err != nil {
		return nil, err
	}

	uu := &models.User{
		ID:       updatedUser.ID,
		Name:     updatedUser.Name,
		Email:    updatedUser.Email,
		UserName: updatedUser.UserName,
	}

	return uu, nil
}

func (s *UserRepository) GetUser(id int64, ctx context.Context) (*models.User, error) {
	user := &User{}

	err := s.db.NewSelect().Model(user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}

	ud := &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		UserName: user.UserName,
	}

	return ud, nil
}

func (s *UserRepository) GetUserByUserName(userName string, ctx context.Context) (*models.User, error) {
	user := &User{}

	err := s.db.NewSelect().Model(user).Where("user_name = ?", userName).Scan(ctx)
	if err != nil {
		return nil, err
	}

	ud := &models.User{
		ID:       user.ID,
		Name:     user.Name,
		Email:    user.Email,
		UserName: user.UserName,
		Password: user.Password,
	}

	return ud, nil
}
