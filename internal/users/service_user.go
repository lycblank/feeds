package users

import "context"

type UserService struct {
	loginRepository LoginRepository
	userRepository UserRepository
}

func NewUserService(loginRepository LoginRepository, userRepository UserRepository) *UserService {
	return &UserService{
		loginRepository:loginRepository,
		userRepository:userRepository,
	}
}

func (us *UserService) Login(ctx context.Context, username string, password string) (*User, error) {
	uid, err := us.loginRepository.Login(ctx, username, password)
	if err != nil {
		return nil, err
	}
	return us.userRepository.GetUser(ctx, uid)
}

func (us *UserService) AddFeed(ctx context.Context, uid int64, feedId int64) error {
	user, err := us.userRepository.GetUser(ctx, uid)
	if err != nil {
		return err
	}
	user.AddFeed(feedId)
	return us.userRepository.SaveUser(ctx, user)
}

func (us *UserService) GetUser(ctx context.Context, uid int64) (*User, error) {
	return us.userRepository.GetUser(ctx, uid)
}




