package user_activation

import (
	"context"

	"necsam/errors"
	"necsam/models"
	"necsam/repos"
)

// UserActivator user activation interface
type UserActivator interface {
	ActivateUser(ctx context.Context, code string) error
}

// UserActivationService UserActivator implementation
type UserActivationService struct {
	userRepo repos.User
	acRepo   repos.ActivationCode
}

// ActivateUser performs user activation
func (svc UserActivationService) ActivateUser(ctx context.Context, code string) error {
	ac, err := svc.acRepo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if ac.IsExpired() {
		return errors.TokenExpiredError
	}

	user, err := svc.userRepo.FindByID(ctx, ac.UserID)
	if err != nil {
		return err
	}

	user.Status = models.StatusUserActive
	if err := svc.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return svc.acRepo.Delete(ctx, ac)
}

// NewUserActivationService UserActivationService constructor
func NewUserActivationService(userRepo repos.User, acRepo repos.ActivationCode) UserActivationService {
	return UserActivationService{userRepo: userRepo, acRepo: acRepo}
}
