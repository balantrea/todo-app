package handler

import (
	"github.com/balantrea/todo-app/internal/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type Handler struct {
	service *service.Service
	logger  zerolog.Logger
}

func NewHandler(service *service.Service, logger zerolog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) InitRouters() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.singUP)
		auth.POST("/sign-in", h.singIn)
	}

	api := router.Group("/api")
	{
		list := api.Group("/list")
		{
			list.POST("/", h.createList)
			list.GET("/", h.getAllList)
			list.GET("/:id", h.getListById)
			list.PUT("/:id", h.getListById)
			list.DELETE("/:id", h.deleteList)

			items := list.Group(":id/items")
			{
				items.POST("/", h.createItems)
				items.GET("/", h.getAllItems)
				items.GET("/:id", h.getItemsById)
				items.PUT("/:id", h.updateItems)
				items.DELETE("/:id", h.deleteItems)
			}
		}
	}

	return router
}
