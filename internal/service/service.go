package service

import (
	"context"
	"errors"
)

// service implements main logic of applicaiton
type service struct {
	// TOOD:
	// Logger
	repository Repository
}

// Configuration settings for service
type Configuration struct {
	// Logger
	Repository Repository
}

// New returns new instance of service
func New(ctx context.Context, conf Configuration) (service, error) {
	if conf.Repository == nil {
		return service{}, errors.New("Repository is not set for service.")
	}

	return service{
		repository: conf.Repository,
	}, nil
}
