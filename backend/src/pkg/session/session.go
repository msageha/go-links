package session

import (
	"github.com/google/uuid"
)

type Session struct {
	Token       string
	AdminUserID uint64
}

const (
	SELECT_OAUTH_ADMIN_SESSION_BY_TOKEN = "SELECT token, user_id FROM session WHERE token = ?"
	INSERT_OAUTH_ADMIN_SESSION          = "INSERT INTO session(token, user_id)"
)

func NewSession(user *u.AdminUser) (*Session, error) {
	token, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	session := &Session{
		Token:       token.String(),
		AdminUserID: user.ID,
	}

	session.TokenType = oauthToken.TokenType
	session.AccessToken = oauthToken.AccessToken
	session.IDToken = oauthToken.IDToken
	session.Scope = oauthToken.Scope
	session.LoginHint = oauthToken.LoginHint
	session.FirstIssuedAt = oauthToken.FirstIssuedAt
	session.ExpiresAt = oauthToken.ExpiresAt
	session.ExpiresIn = oauthToken.ExpiresIn
	session.IdpID = oauthToken.IdpID

	return session, nil
}
