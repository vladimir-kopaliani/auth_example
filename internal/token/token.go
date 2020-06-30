package token

import "time"

// Token contains access, refresh tokens, of user with guid
type Token struct {
	GUID             string    `bson:"guid"`
	Access           string    `bson:"access_token"`
	Refresh          string    `bson:"refresh_token"`
	CreatedAt        time.Time `bson:"created_at"`
	RefreshExpiredAt time.Time `bson:"refresh_expired_at"`
}
