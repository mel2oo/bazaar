package domain

import (
	"bazaar/config"
	"bazaar/internal/domain/browse"
	"bazaar/pkg/logger"
)

type Domain struct {
	*browse.Browse
}

func New(c *config.Config) (*Domain, error) {
	err := logger.Init(c.Server.Name, c.Logger, false, nil)
	if err != nil {
		return nil, err
	}

	db, err := browse.New(c)
	if err != nil {
		return nil, err
	}

	return &Domain{
		Browse: db,
	}, nil
}
