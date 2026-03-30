package handler

import (
	"practice_IstraNet/pkg/service"

	"github.com/gin-gonic/gin"
)

const frontendDir = "../frontend"

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.Static("/style", frontendDir+"/style")
	router.Static("/js", frontendDir+"/js")
	router.Static("/image", frontendDir+"/image")
	router.Static("/components", frontendDir+"/components")

	router.GET("/", func(c *gin.Context) {
		c.File(frontendDir + "/index.html")
	})

	router.GET("/index.html", func(c *gin.Context) {
		c.File(frontendDir + "/index.html")
	})

	router.GET("/auth", func(c *gin.Context) {
		c.File(frontendDir + "/auth.html")
	})

	router.GET("/auth.html", func(c *gin.Context) {
		c.File(frontendDir + "/auth.html")
	})

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	apps := router.Group("/applications", h.userIdentity())
	{
		apps.POST("/create-app", h.CreateApplication)   // role=user
		apps.GET("/get-all-apps", h.GetAllApplications) // role=operator/admin
		apps.GET("/get-apps", h.GetUserApps)            // role=user
		apps.GET("/get-app", h.GetUserApp)              // role=user
		apps.DELETE("/delete-app", h.DeleteUserApp)     // role=user
		apps.PUT("/change-app", h.ChangeUserApp)        // role=user
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

	adminApps := router.Group("/api/admin/applications", h.userIdentity())
	{
		adminApps.GET("/get-apps", h.AdminGetApps)
		adminApps.GET("/get-app", h.AdminGetApp)
		adminApps.PUT("/assign-app", h.AdminAssignApp)
		adminApps.PUT("/change-status", h.AdminChangeStatus)
		adminApps.GET("/get-history", h.AdminGetHistory)
		adminApps.DELETE("/delete-app", h.AdminDeleteApp)
	}

	adminUsers := router.Group("/api/admin/users", h.userIdentity())
	{
		adminUsers.GET("/get-users", h.AdminGetUsers)
		adminUsers.GET("/get-user", h.AdminGetUser)
		adminUsers.PUT("/change-role", h.AdminChangeRole)
		adminUsers.DELETE("/delete-user", h.AdminDeleteUser)
	}

	adminDictionaries := router.Group("/api/admin/dictionaries", h.userIdentity())
	{
		adminDictionaries.POST("/create-status", h.AdminCreateStatus)
		adminDictionaries.PUT("/change-status", h.AdminUpdateStatus)
		adminDictionaries.DELETE("/delete-status", h.AdminDeleteStatus)
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
