package handler

import (
	"net/http"
	"strconv"

	"github.com/balantrea/todo-app/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createItems(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param", h.logger)
		return
	}

	var input model.TodoItem

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	id, err := h.service.TodoItem.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type GetAllItemsResponse struct {
	Data []model.TodoItem `json:"data"`
}

func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param", h.logger)
		return
	}

	items, err := h.service.TodoItem.GetAll(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, GetAllItemsResponse{
		Data: items,
	})
}

func (h *Handler) getItemsById(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param", h.logger)
		return
	}

	item, err := h.service.TodoItem.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *Handler) updateItems(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param", h.logger)
		return
	}

	var input model.UpdateItemInput

	if err = c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	err = h.service.Update(userId, itemId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteItems(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param", h.logger)
		return
	}

	if err := h.service.TodoItem.Delete(userId, itemId); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}
