package users

import "context"

type UserRepository interface {
	GetUser(ctx context.Context, uid int64) (*User, error)
	SaveUser(ctx context.Context, user *User) error
}


