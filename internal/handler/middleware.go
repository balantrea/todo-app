package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header", h.logger)
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header", h.logger)
		return
	}

	userId, err := h.service.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error(), h.logger)
	}

	c.Set(userCtx, userId)
}

func (h *Handler) getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "user is not found", h.logger)
		return 0, fmt.Errorf("user id is not fond")
	}

	userId, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusBadRequest, "user is not found", h.logger)
		return 0, fmt.Errorf("failed to assert user id")
	}

	return userId, nil
}
