package login

import (
	"context"

	"github.com/asaskevich/govalidator"
)

// Form validates login params
type Form struct {
	Email    string `valid:"email" json:"email"`
	Password string `valid:"stringlength(6|40)" json:"password"`
}

// Submit performs validation
func (form Form) Submit(ctx context.Context) (map[string]string, error) {
	_, err := govalidator.ValidateStruct(&form)
	return govalidator.ErrorsByField(err), nil
}

// NewForm LoginForm constructor
func NewForm() Form {
	return Form{}
}
