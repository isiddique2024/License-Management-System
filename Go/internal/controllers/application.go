package controllers

import (
	"backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

// @Summary Create a new application
// @Tags private
// @Description Create a new application for the authenticated user
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started" default(Bearer <token>)
// @Param application body models.CreateApplicationRequest true "Application data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/private/applications [post]
func CreateApplication(ctx *gin.Context, db *gorm.DB) {
	request := ctx.MustGet("request").(*models.CreateApplicationRequest)

	userInfo := ctx.MustGet("userInfo").(map[string]interface{})
	userID := userInfo["sub"].(string)

	appID := uuid.New().String()

	application := models.Application{
		ApplicationID: appID,
		AppName:       request.AppName,
		UserID:        userID,
	}

	if err := db.Create(&application).Error; err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(fasthttp.StatusOK, gin.H{"application": application})
}

// Need to add:

// deleteApplication
// pauseApplication
