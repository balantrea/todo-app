package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Err struct {
	Err string `json:"message"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, err error, logger zerolog.Logger) {
	logger.Info().
		Err(err)

	c.AbortWithStatusJSON(statusCode, Err{err.Error()})
}
