package token

import "golang.org/x/oauth2"

// TokenSource
type TokenSource struct {
	AccessToken string
}

// Token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
