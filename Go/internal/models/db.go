package models

import (
	"gorm.io/gorm"
)

// Application model
type Application struct {
	gorm.Model
	ApplicationID string `gorm:"primaryKey;size:36"`
	AppName       string `gorm:"not null;uniqueIndex:idx_appname_userid"`
	UserID        string `gorm:"not null;size:36;uniqueIndex:idx_appname_userid"`
}

// License model
type License struct {
	gorm.Model
	UserID        string `gorm:"index;size:36"`                                     // UUID is 36 characters
	ApplicationID string `gorm:"size:36;not null;uniqueIndex:idx_application_key"`  // UUID is 36 characters
	Key           string `gorm:"size:100;not null;uniqueIndex:idx_application_key"` // Composite unique index with ApplicationID
	Note          string `gorm:"size:255"`                                          // Limiting note to 255 characters
	CreatedOn     string `gorm:"size:50"`
	Duration      string `gorm:"size:50"`
	GeneratedBy   string `gorm:"size:50"`
	UsedOn        string `gorm:"size:50"`
	ExpiresOn     string `gorm:"size:50"`
	Status        string `gorm:"size:50"`
	IP            string `gorm:"size:45"`  // IPv6 can be up to 45 characters
	HWID          string `gorm:"size:255"` // Can vary
}
