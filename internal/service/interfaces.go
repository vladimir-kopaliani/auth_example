package service

import (
	"context"

	"github.com/vladimir-kopaliani/auth_example/internal/token"
)

// Repository represents database handler
type Repository interface {
	SaveToken(ctx context.Context, tkn *token.Token) error
	GetToken(ctx context.Context, guid, accessToken string) (*token.Token, error)
	ChangeToken(ctx context.Context, guid, oldAccessToken, refreshToken, accessToken string) error
	RemoveToken(ctx context.Context, guid, accessToken string) error
	RemoveAllTokens(ctx context.Context, guid string) error
}
