package controllers

import (
	"log"
	"time"

	"backend/internal/models"
	"backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

// GenerateLicense handles the creation of new licenses for an application.
// @Summary Generate licenses
// @Tags Licenses
// @Description Generate licenses based on provided data
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param application_id path string true "Application ID"
// @Param request body models.LicenseRequest true "License generation data"
// @Success 201 {object} map[string]interface{} "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/private/applications/{application_id}/licenses [post]
func GenerateLicense(ctx *gin.Context, db *gorm.DB) {
	redisClient, ok := ctx.MustGet("redisClient").(*redis.Client)
	if !ok {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to access Redis client",
			"REDIS_CLIENT_ERROR",
			nil,
		))
		return
	}
	request := ctx.MustGet("request").(*models.LicenseRequest)
	applicationID := ctx.Param("application_id")
	userInfo := ctx.MustGet("userInfo").(map[string]interface{})
	userID := userInfo["sub"].(string)
	username := userInfo["preferred_username"].(string)
	currentDateTime := utils.GetCurrentDatetime()

	var application models.Application
	if err := db.Where("application_id = ? AND user_id = ?", applicationID, userID).First(&application).Error; err != nil {
		ctx.JSON(fasthttp.StatusNotFound, utils.NewErrorResponse(
			fasthttp.StatusNotFound,
			"Application not found or does not belong to the user",
			"APPLICATION_NOT_FOUND",
			nil,
		))
		return
	}

	var dbLicenses []models.License
	var licenses []models.LicenseResponse
	for i := 0; i < request.LicenseAmount; i++ {
		key := utils.GenerateLicenseKey(request.Prefix, request.LicenseMask)
		licenseData := models.LicenseResponse{
			Key:         key,
			Note:        request.LicenseNote,
			CreatedOn:   currentDateTime,
			Duration:    utils.FormatDuration(request.LicenseDuration, request.LicenseExpiryUnit),
			GeneratedBy: username,
			UsedOn:      "N/A",
			ExpiresOn:   "N/A",
			Status:      "Not Used",
			IP:          "N/A",
			HWID:        "N/A",
		}

		dbLicenses = append(dbLicenses, models.License{
			UserID:        userID,
			ApplicationID: applicationID,
			Key:           key,
			Note:          request.LicenseNote,
			CreatedOn:     currentDateTime,
			Duration:      utils.FormatDuration(request.LicenseDuration, request.LicenseExpiryUnit),
			GeneratedBy:   username,
			UsedOn:        "N/A",
			ExpiresOn:     "N/A",
			Status:        "Not Used",
			IP:            "N/A",
			HWID:          "N/A",
		})

		licenses = append(licenses, licenseData)
	}

	if err := db.Create(&dbLicenses).Error; err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to create licenses",
			"LICENSE_CREATION_FAILED",
			nil,
		))
		return
	}

	// Invalidate the cache for the application's licenses
	licensesCacheKey := "application:" + applicationID + ":licenses"
	redisClient.Del(ctx, licensesCacheKey)
	log.Printf("Cache invalidated for application %s licenses", applicationID)

	ctx.JSON(fasthttp.StatusCreated, gin.H{"licenses": licenses, "user": userInfo})
}

// RedeemLicense handles the redemption of a license using a license key and HWID.
// @Summary Redeem a license
// @Tags Licenses
// @Description Redeem a license using license key and HWID for a specific application
// @Accept json
// @Produce json
// @Param application_id path string true "Application ID"
// @Param request body models.RedeemLicenseRequest true "License redemption data"
// @Success 200 {object} map[string]string "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 403 {object} map[string]string "Forbidden"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 409 {object} map[string]string "Conflict"
// @Failure 410 {object} map[string]string "Gone"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/public/applications/{application_id}/redeem-license [post]
func RedeemLicense(ctx *gin.Context, db *gorm.DB) {
	redisClient, ok := ctx.MustGet("redisClient").(*redis.Client)
	if !ok {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to access Redis client",
			"REDIS_CLIENT_ERROR",
			nil,
		))
		return
	}

	request := ctx.MustGet("request").(*models.RedeemLicenseRequest)
	applicationID := ctx.Param("application_id")

	var license models.License
	if err := db.Where("key = ? AND application_id = ?", request.Key, applicationID).First(&license).Error; err != nil {
		ctx.JSON(fasthttp.StatusNotFound, utils.NewErrorResponse(
			fasthttp.StatusNotFound,
			"License not found",
			"LICENSE_NOT_FOUND",
			nil,
		))
		return
	}

	if license.Status == "Used" {
		expiresOn, err := time.Parse("2006-01-02 @ 03:04 PM", license.ExpiresOn)
		if err != nil {
			ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
				fasthttp.StatusInternalServerError,
				"Failed to parse expiry date",
				"PARSE_ERROR",
				nil,
			))
			return
		}

		if time.Now().After(expiresOn) {
			ctx.JSON(fasthttp.StatusGone, utils.NewErrorResponse(
				fasthttp.StatusGone,
				"License expired",
				"LICENSE_EXPIRED",
				nil,
			))
			return
		}

		if license.HWID != request.HWID {
			ctx.JSON(fasthttp.StatusConflict, utils.NewErrorResponse(
				fasthttp.StatusConflict,
				"HWID mismatch for used license",
				"HWID_MISMATCH",
				nil,
			))
			return
		}
	} else if license.Status == "Banned" {
		ctx.JSON(fasthttp.StatusForbidden, utils.NewErrorResponse(
			fasthttp.StatusForbidden,
			"License banned",
			"LICENSE_BANNED",
			nil,
		))
		return
	}

	clientIP := utils.GetClientIP(ctx)
	currentTime := time.Now()
	expiresOn := utils.CalculateExpiryDateFromText(license.Duration)

	license.UsedOn = currentTime.Format("2006-01-02 @ 03:04 PM")
	license.Status = "Used"
	license.IP = clientIP
	license.HWID = request.HWID
	license.ExpiresOn = expiresOn

	if err := db.Save(&license).Error; err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to update license",
			"UPDATE_FAILED",
			nil,
		))
		return
	}

	// Invalidate the cache for the application's licenses
	licensesCacheKey := "application:" + applicationID + ":licenses"
	redisClient.Del(ctx, licensesCacheKey)
	log.Printf("Cache invalidated for application %s licenses", applicationID)

	ctx.JSON(fasthttp.StatusOK, gin.H{"message": "Successfully logged in", "expires_on": license.ExpiresOn})
}

// DeleteLicense handles the deletion of a single license.
// @Summary Delete a license
// @Tags Licenses
// @Description Delete a license based on key, application_id, and token
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param application_id path string true "Application ID"
// @Param license_id path string true "License ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/private/applications/{application_id}/licenses/{license_id} [delete]
func DeleteLicense(ctx *gin.Context, db *gorm.DB) {
	redisClient, ok := ctx.MustGet("redisClient").(*redis.Client)
	if !ok {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to access Redis client",
			"REDIS_CLIENT_ERROR",
			nil,
		))
		return
	}

	licenseID := ctx.Param("license_id")
	applicationID := ctx.Param("application_id")

	userInfo := ctx.MustGet("userInfo").(map[string]interface{})
	userID := userInfo["sub"].(string)

	var license models.License
	if err := db.Where("key = ? AND application_id = ? AND user_id = ?", licenseID, applicationID, userID).First(&license).Error; err != nil {
		ctx.JSON(fasthttp.StatusNotFound, utils.NewErrorResponse(
			fasthttp.StatusNotFound,
			"License not found or does not belong to the user",
			"LICENSE_NOT_FOUND",
			nil,
		))
		return
	}

	if err := db.Delete(&license).Error; err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to delete license",
			"DELETE_FAILED",
			nil,
		))
		return
	}

	// Invalidate the cache for the application's licenses
	licensesCacheKey := "application:" + applicationID + ":licenses"
	redisClient.Del(ctx, licensesCacheKey)
	log.Printf("Cache invalidated for application %s licenses after deletion", applicationID)

	ctx.JSON(fasthttp.StatusNoContent, nil)
}

// DeleteLicenses handles the deletion of multiple licenses.
// @Summary Delete multiple licenses
// @Tags Licenses
// @Description Delete multiple licenses based on keys, application_id, and token
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param application_id path string true "Application ID"
// @Param request body models.DeleteLicensesRequest true "Request with license keys to delete"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/private/applications/{application_id}/licenses [delete]
func DeleteLicenses(ctx *gin.Context, db *gorm.DB) {
	redisClient, ok := ctx.MustGet("redisClient").(*redis.Client)
	if !ok {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to access Redis client",
			"REDIS_CLIENT_ERROR",
			nil,
		))
		return
	}

	request := ctx.MustGet("request").(*models.DeleteLicensesRequest)
	applicationID := ctx.Param("application_id")

	userInfo := ctx.MustGet("userInfo").(map[string]interface{})
	userID := userInfo["sub"].(string)

	tx := db.Begin()
	if err := tx.Where("key IN ? AND application_id = ? AND user_id = ?", request.Keys, applicationID, userID).Delete(&models.License{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to delete licenses",
			"DELETE_FAILED",
			nil,
		))
		return
	}

	tx.Commit()

	// Invalidate the cache for the application's licenses
	licensesCacheKey := "application:" + applicationID + ":licenses"
	redisClient.Del(ctx, licensesCacheKey)
	log.Printf("Cache invalidated for application %s licenses after deletion", applicationID)

	ctx.JSON(fasthttp.StatusNoContent, nil)
}

// DeleteAllLicenses handles the deletion of all licenses for an application.
// @Summary Delete all licenses
// @Tags Licenses
// @Description Delete all licenses based on application_id and token
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param application_id path string true "Application ID"
// @Success 204 "No Content"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/private/applications/{application_id}/licenses-all [delete]
func DeleteAllLicenses(ctx *gin.Context, db *gorm.DB) {
	redisClient, ok := ctx.MustGet("redisClient").(*redis.Client)
	if !ok {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to access Redis client",
			"REDIS_CLIENT_ERROR",
			nil,
		))
		return
	}

	applicationID := ctx.Param("application_id")

	userInfo := ctx.MustGet("userInfo").(map[string]interface{})
	userID := userInfo["sub"].(string)

	tx := db.Begin()
	if err := tx.Where("user_id = ? AND application_id = ?", userID, applicationID).Delete(&models.License{}).Error; err != nil {
		tx.Rollback()
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to delete licenses",
			"DELETE_FAILED",
			nil,
		))
		return
	}

	tx.Commit()

	// Invalidate the cache for the application's licenses
	licensesCacheKey := "application:" + applicationID + ":licenses"
	redisClient.Del(ctx, licensesCacheKey)
	log.Printf("Cache invalidated for application %s licenses after deletion", applicationID)

	ctx.JSON(fasthttp.StatusNoContent, nil)
}

// BanLicense handles banning a license.
// @Summary Ban a license
// @Tags Licenses
// @Description Ban a license based on key, application_id, and token
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param application_id path string true "Application ID"
// @Param request body models.BanLicenseRequest true "Request with license key to ban"
// @Success 200 {object} map[string]string "OK"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/private/applications/{application_id}/licenses/{license_id}/ban [patch]
func BanLicense(ctx *gin.Context, db *gorm.DB) {
	redisClient, ok := ctx.MustGet("redisClient").(*redis.Client)
	if !ok {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to access Redis client",
			"REDIS_CLIENT_ERROR",
			nil,
		))
		return
	}

	request := ctx.MustGet("request").(*models.BanLicenseRequest)
	applicationID := ctx.Param("application_id")

	userInfo := ctx.MustGet("userInfo").(map[string]interface{})
	userID := userInfo["sub"].(string)

	var license models.License
	if err := db.Where("key = ? AND application_id = ? AND user_id = ?", request.Key, applicationID, userID).First(&license).Error; err != nil {
		ctx.JSON(fasthttp.StatusNotFound, utils.NewErrorResponse(
			fasthttp.StatusNotFound,
			"License not found",
			"LICENSE_NOT_FOUND",
			nil,
		))
		return
	}

	license.Status = "Banned"
	if err := db.Save(&license).Error; err != nil {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to ban license",
			"BAN_FAILED",
			nil,
		))
		return
	}

	// Invalidate the cache for the application's licenses
	licensesCacheKey := "application:" + applicationID + ":licenses"
	redisClient.Del(ctx, licensesCacheKey)
	log.Printf("Cache invalidated for application %s licenses after deletion", applicationID)

	ctx.JSON(fasthttp.StatusOK, gin.H{"message": "License banned successfully"})
}

// GetData retrieves data based on the token provided.
// @Summary Get data
// @Tags Data
// @Description Get data based on the token
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Success 200 {object} map[string]interface{} "OK"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /api/v1/private/applications/data [get]
func GetData(ctx *gin.Context, db *gorm.DB) {
	redisClient, ok := ctx.MustGet("redisClient").(*redis.Client)
	if !ok {
		ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
			fasthttp.StatusInternalServerError,
			"Failed to access Redis client",
			"REDIS_CLIENT_ERROR",
			nil,
		))
		return
	}

	userInfo := ctx.MustGet("userInfo").(map[string]interface{})
	userID := userInfo["sub"].(string)

	// Cache key for user's applications (now in the case you do any modifications to the Application table directly for the associated userID, then invalidate userID cache)
	applicationsCacheKey := "user:" + userID + ":applications"

	var applications []models.Application

	// Try to get applications from cache
	cachedApplications, err := redisClient.Get(ctx, applicationsCacheKey).Result()
	if err == nil {
		err = json.Unmarshal([]byte(cachedApplications), &applications)
		if err != nil {
			log.Printf("Cache parse error for user %s applications: %v", userID, err)
		} else {
			log.Printf("Cache hit: Retrieved applications from cache for user %s", userID)
		}
	}

	// If cache miss or error, query the database
	if err != nil || applications == nil {
		if err := db.Where("user_id = ?", userID).Find(&applications).Error; err != nil {
			ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
				fasthttp.StatusInternalServerError,
				"Failed to retrieve applications",
				"APPLICATION_RETRIEVAL_FAILED",
				nil,
			))
			return
		}
		// Cache the applications
		applicationsJSON, _ := json.Marshal(applications)
		redisClient.Set(ctx, applicationsCacheKey, applicationsJSON, 10*time.Minute)
		log.Printf("Cache miss: Retrieved applications from database and cached for user %s", userID)
	}

	licensesByApp := make(map[string][]models.LicenseResponse)
	var response []map[string]interface{}

	for _, app := range applications {
		licensesCacheKey := "application:" + app.ApplicationID + ":licenses"
		var licenses []models.License

		// Try to get licenses from cache
		cachedLicenses, err := redisClient.Get(ctx, licensesCacheKey).Result()
		if err == nil {
			err = json.Unmarshal([]byte(cachedLicenses), &licenses)
			if err != nil {
				log.Printf("Cache parse error for licenses of application %s: %v", app.ApplicationID, err)
			} else {
				log.Printf("Cache hit: Retrieved licenses from cache for application %s", app.ApplicationID)
			}
		}

		// If cache miss or error, query the database
		if err != nil || licenses == nil {
			if err := db.Where("application_id = ?", app.ApplicationID).Find(&licenses).Error; err != nil {
				ctx.JSON(fasthttp.StatusInternalServerError, utils.NewErrorResponse(
					fasthttp.StatusInternalServerError,
					"Failed to retrieve licenses",
					"LICENSE_RETRIEVAL_FAILED",
					nil,
				))
				return
			}
			// Cache the licenses
			licensesJSON, _ := json.Marshal(licenses)
			redisClient.Set(ctx, licensesCacheKey, licensesJSON, 10*time.Minute)
			log.Printf("Cache miss: Retrieved licenses from database and cached for application %s", app.ApplicationID)
		}

		for _, license := range licenses {
			licenseResponse := models.LicenseResponse{
				Key:         license.Key,
				Note:        license.Note,
				CreatedOn:   license.CreatedOn,
				Duration:    license.Duration,
				GeneratedBy: license.GeneratedBy,
				UsedOn:      license.UsedOn,
				ExpiresOn:   license.ExpiresOn,
				Status:      license.Status,
				IP:          license.IP,
				HWID:        license.HWID,
			}
			licensesByApp[license.ApplicationID] = append(licensesByApp[license.ApplicationID], licenseResponse)
		}

		appData := map[string]interface{}{
			"application_id": app.ApplicationID,
			"app_name":       app.AppName,
			"created_at":     app.CreatedAt,
			"updated_at":     app.UpdatedAt,
			"licenses":       licensesByApp[app.ApplicationID],
		}
		response = append(response, appData)
	}

	ctx.JSON(fasthttp.StatusOK, gin.H{
		"applications": response,
		"user":         userInfo,
	})
}

// func GetData(ctx *gin.Context, db *gorm.DB) {
// 	userInfo := ctx.MustGet("userInfo").(map[string]interface{})
// 	userID := userInfo["sub"].(string)

// 	type result struct {
// 		applications []models.Application
// 		licenses     []models.License
// 		err          error
// 	}

// 	resultChan := make(chan result)
// 	go func() {
// 		var applications []models.Application
// 		var licenses []models.License

// 		if err := db.Where("user_id = ?", userID).Find(&applications).Error; err != nil {
// 			resultChan <- result{err: err}
// 			return
// 		}

// 		if err := db.Where("user_id = ?", userID).Find(&licenses).Error; err != nil {
// 			resultChan <- result{err: err}
// 			return
// 		}

// 		resultChan <- result{applications: applications, licenses: licenses}
// 	}()

// 	res := <-resultChan
// 	if res.err != nil {
// 		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse(
// 			http.StatusInternalServerError,
// 			"Failed to retrieve data",
// 			"DATA_RETRIEVAL_FAILED",
// 			nil,
// 		))
// 		return
// 	}

// 	licensesByApp := make(map[string][]models.LicenseResponse)
// 	for _, license := range res.licenses {
// 		licenseResponse := models.LicenseResponse{
// 			Key:         license.Key,
// 			Note:        license.Note,
// 			CreatedOn:   license.CreatedOn,
// 			Duration:    license.Duration,
// 			GeneratedBy: license.GeneratedBy,
// 			UsedOn:      license.UsedOn,
// 			ExpiresOn:   license.ExpiresOn,
// 			Status:      license.Status,
// 			IP:          license.IP,
// 			HWID:        license.HWID,
// 		}
// 		licensesByApp[license.ApplicationID] = append(licensesByApp[license.ApplicationID], licenseResponse)
// 	}

// 	var response []map[string]interface{}
// 	for _, app := range res.applications {
// 		appData := map[string]interface{}{
// 			"application_id": app.ApplicationID,
// 			"app_name":       app.AppName,
// 			"created_at":     app.CreatedAt,
// 			"updated_at":     app.UpdatedAt,
// 			"licenses":       licensesByApp[app.ApplicationID],
// 		}
// 		response = append(response, appData)
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{
// 		"applications": response,
// 		"user":         userInfo,
// 	})
// }
