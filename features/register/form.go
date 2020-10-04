package register

import (
	"context"

	"github.com/asaskevich/govalidator"

	"necsam/errors"
	"necsam/repos"
)

type RegisterForm struct {
	Username string     `valid:"stringlength(3|40)~Username must be between 3 and 40 chars,alphanum~Username must be alphanumeric" json:"username"`
	Password string     `valid:"stringlength(6|40)~Password must be between 6 and 40 chars" json:"password"`
	Email    string     `valid:"email~Email is not a valid email" json:"email"`
	repo     repos.User `valid:"-" json:"-"`
}

func (form RegisterForm) Submit(ctx context.Context) (map[string]string, error) {
	_, err := govalidator.ValidateStruct(&form)
	messages := govalidator.ErrorsByField(err)

	if messages["email"] == "" {
		_, err := form.repo.FindByEmail(ctx, form.Email)
		if err == nil {
			messages["email"] = "Email is taken"
			return messages, nil
		}
		if err == errors.RecordNotFound {
			return messages, nil
		}
		if err != nil {
			return messages, err
		}
	}
	return messages, nil
}

func NewForm(repo repos.User) RegisterForm {
	return RegisterForm{repo: repo}
}
