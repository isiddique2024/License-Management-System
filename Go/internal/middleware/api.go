package middleware

import (
	"backend/internal/models"
	"backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// KeycloakAuth checks for user info and validates the access token
func KeycloakAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, err := utils.GetTokenFromHeader(ctx)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(
				http.StatusUnauthorized,
				"Unauthorized access",
				"INVALID_TOKEN",
				nil, // No sensitive details exposed
			))
			ctx.Abort()
			return
		}

		userInfo, err := utils.GetKeycloakUserInfo(accessToken)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse(
				http.StatusUnauthorized,
				"Unauthorized access",
				"INVALID_TOKEN",
				nil, // No sensitive details exposed
			))
			ctx.Abort()
			return
		}

		ctx.Set("userInfo", userInfo)
		ctx.Next()
	}
}

// ParamValidation validates the parameters
func ParamValidation(paramNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, paramName := range paramNames {
			id := c.Param(paramName)
			var err error
			if paramName == "license_id" {
				err = models.ValidateLicenseID(id)
			} else {
				err = models.ValidateUUID(id)
			}
			if err != nil {
				c.JSON(http.StatusBadRequest, utils.NewErrorResponse(
					http.StatusBadRequest,
					"Invalid parameter",
					"INVALID_PARAMETER",
					map[string]string{"parameter": paramName, "error": err.Error()},
				))
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

// JSONValidation validates the JSON request body and validates fields
func JSONValidation(request interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := ctx.ShouldBindJSON(request); err != nil {
			ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(
				http.StatusBadRequest,
				"Invalid request format",
				"INVALID_JSON",
				err.Error(),
			))
			ctx.Abort()
			return
		}

		if validator, ok := request.(models.Validator); ok {
			if err := validator.Validate(); err != nil {
				ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse(
					http.StatusBadRequest,
					"Request validation failed",
					"VALIDATION_ERROR",
					err.Error(),
				))
				ctx.Abort()
				return
			}
		}

		ctx.Set("request", request)
		ctx.Next()
	}
}
