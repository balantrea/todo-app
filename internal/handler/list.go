package handler

import (
	"net/http"
	"strconv"

	"github.com/balantrea/todo-app/internal/model"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	var input model.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	id, err := h.service.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllResponse struct {
	Data []model.TodoList `json:"data"`
}

func (h *Handler) getAllList(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.service.TodoList.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, getAllResponse{
		Data: lists,
	})
}

func (h *Handler) getListById(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param", h.logger)
		return
	}

	lists, err := h.service.TodoList.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, lists)
}

func (h *Handler) updateList(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param", h.logger)
		return
	}

	var input model.UpdateListInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error(), h.logger)
		return
	}

	err = h.service.UpdateList(userId, id, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}

func (h *Handler) deleteList(c *gin.Context) {
	userId, err := h.getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param", h.logger)
		return
	}

	err = h.service.TodoList.DeleteList(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error(), h.logger)
		return
	}

	c.JSON(http.StatusOK, StatusResponse{
		Status: "ok",
	})
}
