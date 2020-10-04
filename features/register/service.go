package register

import (
	"context"
	"time"

	"necsam/models"
	"necsam/repos"
	"necsam/services"
)

type RegisterService interface {
	Register(ctx context.Context, username, email, password string) (models.User, error)
}

type RegisterServiceImpl struct {
	userRepo repos.User
	acRepo   repos.ActivationCode
	uuidgen  func() (string, error)
	crypter  services.Crypter
	mailer   services.ActivationCodeSender
}

func (svc RegisterServiceImpl) Register(ctx context.Context, username, email, password string) (models.User, error) {
	user := models.User{Username: username, Email: email}
	hashedPassword, err := svc.crypter.Hash(password)
	if err != nil {
		return user, err
	}
	user.Password = hashedPassword
	userID, err := svc.uuidgen()
	if err != nil {
		return user, err
	}
	user.UserID = "u-" + userID
	user.Status = models.StatusUserInactive
	if err := svc.userRepo.Create(ctx, user); err != nil {
		return user, err
	}
	code, err := svc.uuidgen()
	if err != nil {
		return user, err
	}
	ac := models.ActivationCode{}
	ac.UserID = user.UserID
	ac.Code = "ac-" + code
	ac.CreatedAt = time.Now()
	if err := svc.acRepo.Create(ctx, ac); err != nil {
		return user, err
	}
	if err := svc.mailer.SendActivationCode(user, ac.Code); err != nil {
		return user, err
	}
	return user, nil
}

func NewRegisterService(repo repos.User, acRepo repos.ActivationCode, crypter services.Crypter, mailer services.ActivationCodeSender, uuidgen func() (string, error)) RegisterServiceImpl {
	return RegisterServiceImpl{userRepo: repo, acRepo: acRepo, crypter: crypter, mailer: mailer, uuidgen: uuidgen}
}
