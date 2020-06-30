package authrepo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	databaseName       = "auth"
	authCollectionName = "users"
)

// Repository with users' credentials
type Repository struct {
	database       *mongo.Database
	client         *mongo.Client
	authCollection *mongo.Collection
}

// Configuration is settings for mongodb connection
type Configuration struct {
	URI      string
	Username string
	Password string
}

// New returns new instance of repository
func New(ctx context.Context, conf *Configuration) (*Repository, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	opts := options.Client()

	if conf != nil {
		if conf.URI == "" {
			opts.ApplyURI("mongodb://localhost:27017")
		} else {
			opts.ApplyURI(conf.URI)
		}

		if conf.Username != "" || conf.Password != "" {
			opts.SetAuth(options.Credential{
				Username: conf.Username,
				Password: conf.Password,
			})
		}
	}

	err := opts.Validate()
	if err != nil {
		return nil, err
	}

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, readpref.PrimaryPreferred())
	if err != nil {
		return nil, err
	}

	log.Println("Connected to Database")

	db := client.Database(databaseName)
	authCollection := db.Collection(authCollectionName)

	r := Repository{
		database:       db,
		client:         client,
		authCollection: authCollection,
	}

	err = r.createIndexes(ctx)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

// Close connetction to database
func (r *Repository) Close(ctx context.Context) error {
	err := r.database.Client().Disconnect(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) createIndexes(ctx context.Context) error {
	indexes := r.authCollection.Indexes()

	_, err := indexes.CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.M{
				"guid": 1,
			},
		},
		{
			Keys: bson.M{
				"access_token": 1,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
