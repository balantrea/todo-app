package handler

import (
	"net/http"
	"strconv"

	"github.com/balantrea/todo-app"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createItems(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input todo.TodoItem

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.service.TodoItem.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type GetAllItemsResponse struct {
	Data []todo.TodoItem `json:"data"`
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	items, err := h.service.TodoItem.GetAll(userId, listId)

	c.JSON(http.StatusOK, GetAllItemsResponse{
		Data: items,
	})
}

func (h *Handler) getItemsById(c *gin.Context) {}

func (h *Handler) updateItems(c *gin.Context) {}

func (h *Handler) deleteItems(c *gin.Context) {}
