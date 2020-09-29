package service

import (
	"context"

	"github.com/beinan/fastid"
	"github.com/hadv/mysql-db-index-benchmark/core/model"
	"github.com/hadv/mysql-db-index-benchmark/core/repo"
	"github.com/pkg/errors"
)

// AccountService account service interface
type AccountService interface {
	Register(ctx context.Context, user *model.User, batchSize int) (*model.User, error)
	BulkInsert(ctx context.Context, users []model.User) error
}

// Account service
type Account struct {
	repo repo.UserRepo
}

// NewAccount create new account service
func NewAccount(repo repo.UserRepo) *Account {
	return &Account{
		repo: repo,
	}
}

// Register create new account
func (a *Account) Register(ctx context.Context, user *model.User) (*model.User, error) {
	users, err := a.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}
	if len(users) > 0 {
		return nil, errors.New("email is already registered")
	}
	if user.Password != user.ConfirmPassword {
		return nil, errors.New("password and confirm password are not matched")
	}
	usr := &model.User{
		ID:        fastid.CommonConfig.GenInt64ID(),
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
		Email:     user.Email,
		Password:  hashAndSalt([]byte(user.Password)),
	}

	user, err = a.repo.Create(ctx, usr)
	if err != nil {
		return nil, errors.Wrap(err, "cannot register new user")
	}

	return user, nil
}

// BulkInsert bulk insert
func (a *Account) BulkInsert(ctx context.Context, users []*model.User, batchSize int) error {
	length := len(users)
	batch := length / batchSize
	for i := 0; i < batch; i++ {
		if err := a.repo.BulkInsert(ctx, users[(i)*batchSize:(i+1)*batchSize]); err != nil {
			return err
		}
	}
	if err := a.repo.BulkInsert(ctx, users[(batch)*batchSize:length]); err != nil {
		return err
	}
	return nil
}
