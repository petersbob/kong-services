package main

import (
	"errors"
	"time"
)

// ServiceTypeCode is unique code identifying the different types of service a user can create and manage
type ServiceTypeCode uint

const (
	ServiceTypeCodeDatabase           ServiceTypeCode = 1
	ServiceTypeCodeReporting          ServiceTypeCode = 2
	ServiceTypeCodeCurrencyConversion ServiceTypeCode = 3
	ServiceTypeCodeTranslation        ServiceTypeCode = 4
	ServiceTypeCodeNotifications      ServiceTypeCode = 5
)

// CurrentServiceTypes is a map of all of the current service types that are available to create.
var currentServiceTypes = map[ServiceTypeCode]ServiceType{
	ServiceTypeCodeDatabase: {
		TypeCode:          ServiceTypeCodeDatabase,
		Name:              "Database Service",
		Description:       "A service for running databases",
		VersionsAvailable: []uint{1, 2, 3},
	},
	ServiceTypeCodeReporting: {
		TypeCode:          ServiceTypeCodeReporting,
		Name:              "Reporting Service",
		Description:       "A service for running reports",
		VersionsAvailable: []uint{55, 7},
	},
	ServiceTypeCodeCurrencyConversion: {
		TypeCode:          ServiceTypeCodeCurrencyConversion,
		Name:              "Currency Conversion",
		Description:       "A service for doing currency convertions",
		VersionsAvailable: []uint{1, 2, 3, 4, 5, 6, 7, 8},
	},
	ServiceTypeCodeTranslation: {
		TypeCode:          ServiceTypeCodeTranslation,
		Name:              "Translation service",
		Description:       "A service for doing language translations",
		VersionsAvailable: []uint{12, 14},
	},
	ServiceTypeCodeNotifications: {
		TypeCode:          ServiceTypeCodeNotifications,
		Name:              "Notifications service",
		Description:       "A service sending notifications",
		VersionsAvailable: []uint{1},
	},
}

// errorNotFound is a general purpose error for missing resources
var errorNotFound = errors.New("not found")

// Service an implementation of a ServiceType
type Service struct {
	ServiceType
	VersionsInUse []uint `json:"versions_in_use"`
}

// ServiceType describes a service type a user can create and manage
type ServiceType struct {
	TypeCode          ServiceTypeCode `json:"type_code"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	VersionsAvailable []uint          `json:"versions_available"`
}

// ServiceType describes a specific version of a ServiceType a user has deployed
type InstalledServiceVersion struct {
	ServiceType   ServiceTypeCode
	VersionNumber uint
	CreatedAt     time.Time
}
