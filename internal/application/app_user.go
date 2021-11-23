package application

import "context"

type UserApp struct {

}

func (ua *UserApp) Login(ctx context.Context, username string, password string) error {
	return nil
}