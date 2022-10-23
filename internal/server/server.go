package server

import (
	"github.com/gin-gonic/gin"
)

func New() (*gin.Engine, error) {
	router := gin.New()
	router.Use(gin.Recovery())

	return router, nil
}
