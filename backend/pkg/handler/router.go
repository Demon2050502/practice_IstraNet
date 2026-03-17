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

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	apps := router.Group("/applications", h.userIdentity())
	{
		apps.POST("/create-app", h.CreateApplication)       // role=user
		apps.GET("/get-all-apps", h.GetAllApplications)     // role=operator/admin
		apps.GET("/get-apps", h.GetUserApps)                // role=user
		apps.GET("/get-app", h.GetUserApp)                  // role=user
		apps.DELETE("/delete-app", h.DeleteUserApp)         // role=user
		apps.PUT("/change-app", h.ChangeUserApp)            // role=user
	}

	operatorApps := router.Group("/api/operator/applications", h.userIdentity())
	{
		operatorApps.GET("/get-apps", h.OperatorGetApps)
		operatorApps.GET("/get-app", h.OperatorGetApp)
		operatorApps.PUT("/take-app", h.OperatorTakeApp)
		operatorApps.PUT("/change-status", h.OperatorChangeStatus)
		operatorApps.GET("/get-history", h.OperatorGetHistory)
		operatorApps.PUT("/close-app", h.OperatorCloseApp)
	}

	test := router.Group("/test")
	{
		test.POST("/status", response)
	}

	return router
}

func response(c *gin.Context) {
	c.JSON(200, gin.H{"ok": true})
}

