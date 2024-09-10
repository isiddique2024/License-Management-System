package utils

import (
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GetAdminAccessToken() (string, error) {
	clientID := os.Getenv("ADMIN_CLIENT_ID")
	clientSecret := os.Getenv("CLIENT_SECRET")
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	realm := os.Getenv("REALM")

	urltoPost := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakURL, realm)
	args := fasthttp.AcquireArgs()
	defer fasthttp.ReleaseArgs(args)

	args.Set("client_id", clientID)
	args.Set("client_secret", clientSecret)
	args.Set("grant_type", "client_credentials")

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(urltoPost)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBodyString(args.String())

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil || resp.StatusCode() != fasthttp.StatusOK {
		return "", fmt.Errorf("failed to get access token")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &result); err != nil {
		return "", err
	}

	token, ok := result["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("access token not found in response")
	}

	return token, nil
}

func GetKeycloakUserInfo(token string) (map[string]interface{}, error) {
	keycloakURL := os.Getenv("KEYCLOAK_URL")
	realm := os.Getenv("REALM")

	url := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", keycloakURL, realm)
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		return nil, fmt.Errorf("invalid token")
	}

	var userInfo map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &userInfo); err != nil {
		return nil, err
	}

	return userInfo, nil
}
