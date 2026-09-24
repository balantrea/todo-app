package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type err struct {
	Message string `json:"message"`
}

type StatusResponse struct {
	Status string `json:"status"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string, logger zerolog.Logger) {
	logger.Info().
		Msg(message)

	c.AbortWithStatusJSON(statusCode, err{message})
}
