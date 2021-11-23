package users

import "context"

type LoginRepository interface {
	Login(ctx context.Context, loginName,password string) (uid int64, err error)
}
