package authrepo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/vladimir-kopaliani/auth_example/internal/token"
)

// SaveToken saves user's refresh token and its guid in auth collecton
func (r *Repository) SaveToken(ctx context.Context, tkn *token.Token) error {
	_, err := r.authCollection.InsertOne(ctx, tkn)
	if err != nil {
		return err
	}

	return nil
}

// GetToken returns whole token by refresh and access token
func (r *Repository) GetToken(ctx context.Context, guid, accessToken string) (*token.Token, error) {
	var tkn token.Token

	session, err := r.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		return nil, err
	}

	err = mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		err := r.authCollection.FindOne(ctx, bson.M{
			"guid":         guid,
			"access_token": accessToken,
		}).Decode(&tkn)
		if err != nil {
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &tkn, nil
}

// ChangeToken changes token by ids guid
func (r *Repository) ChangeToken(ctx context.Context, guid, oldAccessToken, refreshToken, accessToken string) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		_, err := r.authCollection.UpdateOne(ctx,
			bson.M{
				"guid":         guid,
				"access_token": oldAccessToken,
			},
			bson.M{
				"$set": bson.M{
					"access_token":  accessToken,
					"refresh_token": refreshToken,
					"created_at":    time.Now().UTC(),
				},
			},
		)
		if err != nil {
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// RemoveToken removes user's token by guid and access token
func (r *Repository) RemoveToken(ctx context.Context, guid, accessToken string) error {
	session, err := r.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	err = mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		_, err := r.authCollection.DeleteOne(ctx, bson.M{
			"guid":         guid,
			"access_token": accessToken,
		})
		if err != nil {
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// RemoveAllTokens remove all user's tokens by guid
func (r *Repository) RemoveAllTokens(ctx context.Context, guid string) error {
	session, err := r.client.StartSession()
	if err != nil {
		log.Println(err)
		return err
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		log.Println(err)
		return err
	}

	err = mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		_, err := r.authCollection.DeleteMany(ctx, bson.M{
			"guid": guid,
		})
		if err != nil {
			log.Println(err)
			return err
		}

		if err = session.CommitTransaction(ctx); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
