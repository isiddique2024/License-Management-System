package controllers

import (
	"backend/internal/middleware"
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

// Global Register Routes
func RegisterRoutes(route *gin.Engine, db *gorm.DB) {
	api := route.Group("/api/v1")
	{
		registerDevRoutes(api)
		registerPrivateRoutes(api, db)
		registerPublicRoutes(api, db)
	}

	// Health check for container
	route.GET("/health", func(c *gin.Context) {
		c.JSON(fasthttp.StatusOK, gin.H{
			"DN": "healthy",
		})
	})
}

// Register Developer Routes
func registerDevRoutes(api *gin.RouterGroup) {
	dev := api.Group("/dev")
	{
		dev.POST("/register", RegisterUser)
		dev.POST("/login", LoginUser)
		dev.POST("/verify-keycloak-token", VerifyKeycloakToken)
	}
}

// Register Private Routes (Keycloak Access Token)
func registerPrivateRoutes(api *gin.RouterGroup, db *gorm.DB) {
	private := api.Group("/private")
	private.Use(middleware.KeycloakAuth()) // Use combined middleware for Keycloak auth and user info check
	{
		private.POST("/applications", middleware.JSONValidation(&models.CreateApplicationRequest{}), func(c *gin.Context) { CreateApplication(c, db) })
		private.POST("/applications/:application_id/licenses", middleware.ParamValidation("application_id"), middleware.JSONValidation(&models.LicenseRequest{}), func(c *gin.Context) { GenerateLicense(c, db) })
		private.DELETE("/applications/:application_id/licenses/:license_id", middleware.ParamValidation("application_id", "license_id"), func(c *gin.Context) { DeleteLicense(c, db) })
		private.DELETE("/applications/:application_id/licenses", middleware.ParamValidation("application_id"), middleware.JSONValidation(&models.DeleteLicensesRequest{}), func(c *gin.Context) { DeleteLicenses(c, db) })
		private.DELETE("/applications/:application_id/licenses-all", middleware.ParamValidation("application_id"), func(c *gin.Context) { DeleteAllLicenses(c, db) })
		private.PATCH("/applications/:application_id/licenses/:license_id/ban", middleware.ParamValidation("application_id", "license_id"), middleware.JSONValidation(&models.BanLicenseRequest{}), func(c *gin.Context) { BanLicense(c, db) })
		private.GET("/applications/data", func(c *gin.Context) { GetData(c, db) })
	}
}

// Register Public Routes (For customers)
func registerPublicRoutes(api *gin.RouterGroup, db *gorm.DB) {
	public := api.Group("/public")
	{
		public.POST("/applications/:application_id/redeem-license", middleware.ParamValidation("application_id"), middleware.JSONValidation(&models.RedeemLicenseRequest{}), func(c *gin.Context) {
			RedeemLicense(c, db)
		})
	}
}
