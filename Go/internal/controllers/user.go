package controllers

import (
	"fmt"
	"os"

	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var (
	json        = jsoniter.ConfigCompatibleWithStandardLibrary
	keycloakURL = os.Getenv("KEYCLOAK_URL")
)

// @Summary Register a new user
// @Tags dev
// @Description Register a new user with email and password
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/dev/register [post]
func RegisterUser(ctx *gin.Context) {
	var user models.UserRequest
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(fasthttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := utils.GetAdminAccessToken()
	if err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
		return
	}

	url := fmt.Sprintf("%s/admin/realms/demo/users", keycloakURL)
	payload := map[string]interface{}{
		"username":  user.Email,
		"enabled":   true,
		"email":     user.Email,
		"firstName": user.Email,
		"lastName":  user.Email,
		"credentials": []map[string]interface{}{
			{
				"type":      "password",
				"value":     user.Password,
				"temporary": false,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, gin.H{"error": "Failed to marshal payload"})
		return
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.SetBody(jsonData)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil || resp.StatusCode() != fasthttp.StatusCreated {
		ctx.JSON(fasthttp.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	ctx.JSON(fasthttp.StatusOK, result)
}

// @Summary Login a user
// @Tags dev
// @Description Login a user with username and password
// @Accept json
// @Produce json
// @Param login body models.LoginRequest true "Login data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /api/v1/dev/login [post]
func LoginUser(ctx *gin.Context) {
	var login models.LoginRequest
	if err := ctx.ShouldBindJSON(&login); err != nil {
		ctx.JSON(fasthttp.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Set("client_id", "real-client") // add to env later
	args.Set("grant_type", "password")
	args.Set("username", login.Username)
	args.Set("password", login.Password)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%s/realms/demo/protocol/openid-connect/token", keycloakURL))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType("application/x-www-form-urlencoded")
	args.WriteTo(req.BodyWriter())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil || resp.StatusCode() != fasthttp.StatusOK {
		ctx.JSON(fasthttp.StatusUnauthorized, gin.H{"error": "Login failed"})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
		return
	}

	ctx.JSON(fasthttp.StatusOK, result)
}
