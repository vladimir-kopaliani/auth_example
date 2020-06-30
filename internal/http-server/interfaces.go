package serverhttp

import "context"

// Service caontain main logic of applicaiotn
type Service interface {
	AuthorizeUser(ctx context.Context, guid, accessToken string) (refreshToken string, err error)
	RefreshUserToken(ctx context.Context, guid, oldAccessToken, oldRefreshToken, newAccessToken string) (string, error)
	RemoveRefreshUserToken(ctx context.Context, guid, accessToken, refreshToken string) error
	RemoveAllTokens(ctx context.Context, guid string) error
}
