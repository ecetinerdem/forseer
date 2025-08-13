package database

import (
	"context"
	"fmt"
	"time"

	"github.com/ecetinerdem/forseer/types"
)

type UserRepo interface {
	CreateUser(context.Context, *types.User) (*types.User, error)
	GetUserById(context.Context, string) (*types.User, error)
	GetUserByEmail(context.Context, string) (*types.User, error)
	UpdateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error // Fixed return type
}

func (db *DB) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	query := `
		INSERT INTO users(name, email, password_hashed)
		VALUES($1, $2, $3)
		RETURNING id, subscription, register_date, last_login, is_admin, is_paid
	`

	var userID string
	var userSubscription types.SubscriptionType
	var userRegisterDate time.Time
	var userLastLogin time.Time
	var userIsAdmin bool
	var userIsPaid bool

	err := db.QueryRowContext(
		ctx, query, user.Name, user.Email, user.PasswordHashed,
	).Scan(
		&userID, &userSubscription, &userRegisterDate, &userLastLogin, &userIsAdmin, &userIsPaid,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = userID
	user.Subscription = userSubscription
	user.RegisterDate = userRegisterDate
	user.LastLogin = userLastLogin
	user.IsAdmin = userIsAdmin
	user.IsPaid = userIsPaid

	return user, nil
}

func (db *DB) GetUserById(ctx context.Context, id string) (*types.User, error) {
	// TODO: Implement
	return nil, nil
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	// TODO: Implement
	return nil, nil
}

func (db *DB) UpdateUser(ctx context.Context, user *types.User) (*types.User, error) {
	// TODO: Implement
	return nil, nil
}

func (db *DB) DeleteUser(ctx context.Context, id string) error {
	// TODO: Implement
	return nil
}
