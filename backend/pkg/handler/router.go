package handler

import (
	"practice_IstraNet/pkg/service"

	"github.com/gin-gonic/gin"
)


type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}


func (h *Handler)InitRoutes() *gin.Engine {
	router := gin.New()

	// auth := router.Group("/auth")
	// {
	// 	auth.POST("/sign-up", h.MainHandler.Authorization.SignUp)
	// 	auth.POST("/sign-in", h.MainHandler.Authorization.SignIn)
	// }

	test := router.Group("/test")
	{
		test.POST("/status", responce)
	}

	return router
}

func responce(c *gin.Context) {
	c.JSON(200, gin.H{"ok": true})
}