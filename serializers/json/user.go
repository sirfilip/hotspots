package json

import "necsam/models"

// User json serializer
type User struct {
	UserID   string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Populate populates user serializer fields from user model
func (u *User) Populate(user models.User) {
	u.UserID = user.UserID
	u.Email = user.Email
	u.Username = user.Username
}
