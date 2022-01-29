package main

type ServiceTypeCode uint

const (
	ServiceTypeCodeDatabase           ServiceTypeCode = 1
	ServiceTypeCodeReporting          ServiceTypeCode = 2
	ServiceTypeCodeCurrencyConversion ServiceTypeCode = 3
	ServiceTypeCodeTranslation        ServiceTypeCode = 4
	ServiceTypeCodeNotifications      ServiceTypeCode = 5
)

type ServiceType struct {
	TypeCode          ServiceTypeCode
	Name              string
	Description       string
	VersionsAvailable []uint
}

var currentServices = map[ServiceTypeCode]ServiceType{
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

func getServiceType(typeCode ServiceTypeCode) ServiceType {
	return currentServices[typeCode]
}
