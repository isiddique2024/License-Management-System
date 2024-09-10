package controllers

import (
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
)

// @Summary Verify Keycloak Token
// @Tags dev
// @Description Verify the provided token in the Authorization header
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/dev/verify-keycloak-token [post]
func VerifyKeycloakToken(ctx *gin.Context) {
	token, err := utils.GetTokenFromHeader(ctx)
	if err != nil {
		ctx.JSON(fasthttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInfo, err := utils.GetKeycloakUserInfo(token)
	if err != nil {
		ctx.JSON(fasthttp.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	ctx.JSON(fasthttp.StatusOK, gin.H{"user_info": userInfo})
}
