package login

import (
	"context"

	"necsam/errors"
	"necsam/models"
	"necsam/repos"
	"necsam/services"
)

// Service login interface
type Service interface {
	Login(ctx context.Context, email, password string) (models.Token, error)
}

type LoginService struct {
	userRepo repos.User
	crypter  services.Crypter
}

func (svc LoginService) Login(ctx context.Context, email, password string) (models.Token, error) {
	token := models.Token{}
	user, err := svc.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return token, err
	}

	if user.Status != models.StatusUserActive {
		return token, errors.RecordNotFound
	}

	if err := svc.crypter.Check(user.Password, password); err != nil {
		return token, errors.RecordNotFound
	}

	accessToken, err := services.NewAccessTokenFor(user)

	if err != nil {
		return token, err
	}
	token.AccessToken = accessToken
	return token, nil
}

func NewService(userRepo repos.User, crypter services.Crypter) LoginService {
	return LoginService{userRepo: userRepo, crypter: crypter}
}
