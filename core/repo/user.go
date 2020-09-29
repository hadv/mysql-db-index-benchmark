package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hadv/mysql-db-index-benchmark/core/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// UserRepo User repository interface
type UserRepo interface {
	BulkInsert(ctx context.Context, users []*model.User) error
	Create(ctx context.Context, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, email string) ([]model.User, error)
}

// User repository
type User struct {
	db *sqlx.DB
}

// NewUser create new user repository
func NewUser(db *sqlx.DB) *User {
	return &User{
		db: db,
	}
}

// BulkInsert bulk insert user
func (u *User) BulkInsert(ctx context.Context, users []*model.User) error {
	length := len(users)
	if length <= 0 {
		return nil
	}
	valueStrings := make([]string, 0, length)
	valueArgs := make([]interface{}, 0, length*6)
	for _, user := range users {
		valueStrings = append(valueStrings, "(?, ?, ?, ?, ?, ?)")
		valueArgs = append(valueArgs, user.ID)
		valueArgs = append(valueArgs, user.Firstname)
		valueArgs = append(valueArgs, user.Lastname)
		valueArgs = append(valueArgs, user.Email)
		valueArgs = append(valueArgs, user.Password)
		valueArgs = append(valueArgs, user.Token)

	}
	stmt := fmt.Sprintf("INSERT INTO `users` (`id`, `firstname`, `lastname`, `email`, `password`, `token`) VALUES %s",
		strings.Join(valueStrings, ","))
	if _, err := u.db.ExecContext(ctx, stmt, valueArgs...); err != nil {
		return err
	}
	return nil
}

// Create insert new user into db
func (u *User) Create(ctx context.Context, user *model.User) (*model.User, error) {
	_, err := u.db.ExecContext(ctx, "INSERT INTO `users` (`id`, `firstname`, `lastname`, `email`, `password`, `token`) VALUES(?, ?, ?, ?, ?, ?)",
		user.ID, user.Firstname, user.Lastname, user.Email, user.Password, user.Token)
	if err != nil {
		return nil, errors.Wrap(err, "cannot insert new user")
	}
	var usr model.User
	err = u.db.GetContext(ctx, &usr, "SELECT * FROM `users` WHERE `id` = ?", user.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find user with id=%s", user.ID)
	}

	return nil, nil
}

// FindByEmail list users by email
func (u *User) FindByEmail(ctx context.Context, email string) ([]model.User, error) {
	var users []model.User
	err := u.db.SelectContext(ctx, &users, "SELECT * FROM `users` WHERE `email` = ?", email)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find user with email=%s", email)
	}
	return users, nil
}
