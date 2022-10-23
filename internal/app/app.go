package app

import (
	"bazaar/config"
	"bazaar/internal/domain"
	"bazaar/internal/router"
	"bazaar/internal/server"
)

func Run(c *config.Config) error {
	srv, err := server.New()
	if err != nil {
		return err
	}

	domain, err := domain.New(c)
	if err != nil {
		return err
	}

	if err := router.New(srv, domain); err != nil {
		return err
	}

	return srv.Run(c.Server.Address)
}
