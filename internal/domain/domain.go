package domain

import (
	"bazaar/config"
	"bazaar/internal/domain/browse"
)

type Domain struct {
	*browse.Browse
}

func New(c *config.Config) (*Domain, error) {
	db, err := browse.New(c)
	if err != nil {
		return nil, err
	}

	return &Domain{
		Browse: db,
	}, nil
}
