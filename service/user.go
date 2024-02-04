package service

import (
	"context"
	"signer/config"
	"signer/repo"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, username string, password string) error
	CheckUser(ctx context.Context, username string, password string) error
}

type userService struct {
	config *config.Config
	ur     *repo.UserRepo
}

func NewUserService(cfg *config.Config, ur *repo.UserRepo) UserService {
	return &userService{config: cfg, ur: ur}
}

func (us *userService) RegisterUser(ctx context.Context, username string, password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	return us.ur.CreateNewUser(username, string(bytes))
}

func (us *userService) CheckUser(ctx context.Context, username string, password string) error {
	pass, err := us.ur.GetUserPassword(username)
	if err != nil {
		return err
	}
	return bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
}
