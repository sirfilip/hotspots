package json

import "necsam/models"

// Token json serializer
type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// Populate populates token serializer fields from token
func (t *Token) Populate(token models.Token) {
	t.AccessToken = token.AccessToken
	t.RefreshToken = token.RefreshToken
}
