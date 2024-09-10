package models

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/microcosm-cc/bluemonday"
)

// Interface for JSON request body validation
type Validator interface {
	Validate() error
}

// sanitizeInput sanitizes a string using bluemonday's StrictPolicy
func sanitizeInput(input string) string {
	policy := bluemonday.StrictPolicy()
	return policy.Sanitize(input)
}

// CreateApplicationRequest is the JSON request body for creating an application
type CreateApplicationRequest struct {
	AppName string `json:"appName" binding:"required"`
}

// Input validation method for CreateApplicationRequest
func (createApplicationRequest *CreateApplicationRequest) Validate() error {
	createApplicationRequest.AppName = sanitizeInput(createApplicationRequest.AppName)

	return validation.ValidateStruct(createApplicationRequest,
		validation.Field(&createApplicationRequest.AppName, validation.RuneLength(0, 25)),
	)
}

// LicenseRequest is the JSON request body for creating a license
type LicenseRequest struct {
	LicenseAmount     int    `json:"license_amount"`
	LicenseMask       string `json:"license_mask"`
	Prefix            string `json:"prefix"`
	LicenseNote       string `json:"license_note"`
	LicenseExpiryUnit string `json:"license_expiry_unit"`
	LicenseDuration   int    `json:"license_duration"`
}

// Input validation method for LicenseRequest
func (licenseRequest *LicenseRequest) Validate() error {
	// add more sanitization here
	licenseRequest.Prefix = sanitizeInput(licenseRequest.Prefix)
	licenseRequest.LicenseNote = sanitizeInput(licenseRequest.LicenseNote)
	licenseRequest.LicenseExpiryUnit = sanitizeInput(licenseRequest.LicenseExpiryUnit)

	return validation.ValidateStruct(licenseRequest,
		validation.Field(&licenseRequest.LicenseAmount, validation.Required, validation.Min(1), validation.Max(100)),
		validation.Field(&licenseRequest.LicenseMask, validation.Required, validation.Match(regexp.MustCompile(`^([X]+(-[X]+)*)?$`))),
		validation.Field(&licenseRequest.Prefix, validation.Required, validation.Length(1, 25), validation.Match(regexp.MustCompile(`^[A-Za-z0-9]+$`))),
		validation.Field(&licenseRequest.LicenseNote, validation.RuneLength(0, 255)),
		validation.Field(&licenseRequest.LicenseExpiryUnit, validation.Required, validation.In("Day", "Days", "Week", "Weeks", "Month", "Months", "Year", "Years")),
		validation.Field(&licenseRequest.LicenseDuration, validation.Required, validation.Min(1), validation.Max(10)),
	)
}

// RedeemLicenseRequest is the JSON request body for redeeming a license
type RedeemLicenseRequest struct {
	Key  string `json:"key" binding:"required"`
	HWID string `json:"hwid" binding:"required"`
}

// Input validation method for RedeemLicenseRequest
func (redeemLicenseRequest *RedeemLicenseRequest) Validate() error {
	redeemLicenseRequest.Key = sanitizeInput(redeemLicenseRequest.Key)
	redeemLicenseRequest.HWID = sanitizeInput(redeemLicenseRequest.HWID)

	return validation.ValidateStruct(redeemLicenseRequest,
		validation.Field(&redeemLicenseRequest.Key, validation.Required, validation.Length(1, 100), validation.Match(regexp.MustCompile(`^[A-Za-z0-9-]+$`))),
		validation.Field(&redeemLicenseRequest.HWID, validation.Required, validation.Length(1, 255)),
	)
}

// DeleteLicensesRequest is the JSON request body for deleting multiple licenses
type DeleteLicensesRequest struct {
	Keys []string `json:"keys" binding:"required"`
}

// Input validation method for DeleteLicensesRequest
func (deleteLicensesRequest *DeleteLicensesRequest) Validate() error {
	for i, key := range deleteLicensesRequest.Keys {
		deleteLicensesRequest.Keys[i] = sanitizeInput(key)
	}

	return validation.ValidateStruct(deleteLicensesRequest,
		validation.Field(&deleteLicensesRequest.Keys,
			validation.Required,
			validation.Each(
				validation.Required,
				validation.Length(1, 100),
				validation.Match(regexp.MustCompile(`^[A-Za-z0-9-]+$`)),
			),
		),
	)
}

// BanLicenseRequest is the JSON request body for banning a license
type BanLicenseRequest struct {
	Key string `json:"key" binding:"required"`
}

// Input validation method for BanLicenseRequest
func (banLicenseRequest *BanLicenseRequest) Validate() error {
	banLicenseRequest.Key = sanitizeInput(banLicenseRequest.Key)

	return validation.ValidateStruct(banLicenseRequest,
		validation.Field(&banLicenseRequest.Key, validation.Required, validation.Length(1, 100), validation.Match(regexp.MustCompile(`^[A-Za-z0-9-]+$`))),
	)
}

// ValidateUUID validates if a given string is a valid UUID with a length of 36.
func ValidateUUID(uuid string) error {
	sanitizedUUID := sanitizeInput(uuid)

	return validation.Validate(sanitizedUUID,
		validation.Required,
		validation.Length(36, 36),
		validation.Match(regexp.MustCompile("^[a-fA-F0-9-]{36}$")),
	)
}

// ValidateLicenseID validates if a given string is a valid License ID with a max length of 100.
func ValidateLicenseID(licenseID string) error {
	sanitizedLicenseID := sanitizeInput(licenseID)

	return validation.Validate(sanitizedLicenseID,
		validation.Required,
		validation.Length(1, 100),
		validation.Match(regexp.MustCompile("^[A-Za-z0-9-]+$")),
	)
}

// dev request below

type UserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
